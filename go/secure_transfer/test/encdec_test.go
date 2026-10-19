package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mats0319/secure_transfer/internal"
	"github.com/mats0319/secure_transfer/utils"
)

// TestRoundTrip 覆盖两种加密方式在各类尺寸下的加解密往返与密文长度。
//
// 尺寸选取覆盖了历史上出过问题的地方：
//   - 0 字节（空明文，流式会产出一个空末帧）
//   - 1 字节 / 不足一帧（末帧不足 FrameSize）
//   - FrameSize 的整数倍与非整数倍（末帧是否为满帧）
//   - OnceEncMaxSize 前后（一次性加密与流式加密的选择边界）
func TestRoundTrip(t *testing.T) {
	sizes := []int64{
		0, 1, 25, 1024, 500 << 10,
		1<<20 - 1, 1 << 20, 1<<20 + 1,
		3<<20 + 12345,
		5 << 20,
		utils.OnceEncMaxSize - 1,
		utils.OnceEncMaxSize,
		utils.OnceEncMaxSize + 1,
		12<<20 + 1,
	}

	for _, mode := range []string{ModeOnce, ModeStream} {
		for _, size := range sizes {
			if mode == ModeOnce && size > utils.OnceEncMaxSize {
				// 超过阈值时生产路径不会选择一次性加密，单独构造这种密文没有意义
				continue
			}

			t.Run(fmt.Sprintf("%s/%s", mode, humanSize(size)), func(t *testing.T) {
				dir := t.TempDir()
				plain := filepath.Join(dir, "plain.bin")
				cipher := filepath.Join(dir, "cipher.bin")
				out := filepath.Join(dir, "plain.out")

				writeRandomFile(t, plain, size)

				priv, pub := newKeyPair(t)
				encryptFile(t, mode, pub, plain, cipher)

				if got := modeOf(t, cipher); got != mode {
					t.Fatalf("cipher header records mode %q, want %q", got, mode)
				}
				requireCipherSize(t, mode, cipher, size)

				requireNoError(t, tryDecrypt(priv, cipher, out), "decrypt")
				requireSameContent(t, plain, out)
			})
		}
	}
}

// TestEncryptMethodSelection 走生产入口，验证按文件大小选择加密方式，
// 并确认边界尺寸选出的方式确实能完整往返。
func TestEncryptMethodSelection(t *testing.T) {
	cases := []struct {
		size int64
		want string
	}{
		{utils.OnceEncMaxSize - 1, ModeOnce},
		{utils.OnceEncMaxSize, ModeOnce}, // 恰好等于阈值仍走一次性加密
		{utils.OnceEncMaxSize + 1, ModeStream},
	}

	for _, c := range cases {
		t.Run(humanSize(c.size), func(t *testing.T) {
			inWorkDir(t)

			if err := internal.GenerateKeyPair(true); err != nil {
				t.Fatalf("generate key pair: %v", err)
			}
			writeRandomFile(t, "message.txt", c.size)

			if err := internal.Encrypt(); err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			if got := modeOf(t, "CIPHER.TXT"); got != c.want {
				t.Fatalf("size %d selected %q, want %q", c.size, got, c.want)
			}

			if err := internal.Decrypt(); err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			requireSameContent(t, "message.txt", "message_decrypted.txt")
		})
	}
}

// TestProductionFileNames 验证生产流程根据输入扩展名派生的输出文件名。
func TestProductionFileNames(t *testing.T) {
	cases := []struct {
		plain     string
		cipher    string
		decrypted string
	}{
		{"message.txt", "CIPHER.TXT", "message_decrypted.txt"},
		{"message.md", "CIPHER.MD", "message_decrypted.md"},
		{"message.TXT", "CIPHER.TXT", "message_decrypted.txt"},
	}

	for _, c := range cases {
		t.Run(c.plain, func(t *testing.T) {
			inWorkDir(t)

			if err := internal.GenerateKeyPair(true); err != nil {
				t.Fatalf("generate key pair: %v", err)
			}
			writeRandomFile(t, c.plain, 1234)

			if err := internal.Encrypt(); err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			if got := fileSize(t, c.cipher); got != int64(1234)+cipherHeaderLen+utils.GCMTagLength {
				t.Fatalf("%s size = %d, want %d", c.cipher, got, int64(1234)+cipherHeaderLen+utils.GCMTagLength)
			}

			if err := internal.Decrypt(); err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			requireSameContent(t, c.plain, c.decrypted)
		})
	}
}

// TestReuseEncryptorInstance 验证同一个 Encryptor 实例被复用两次时，
// 每次都会生成新的临时密钥对 / salt / nonce：
// 两份密文必须都能正确解密，且密文（进而临时密钥）不能相同。
func TestReuseEncryptorInstance(t *testing.T) {
	for _, mode := range []string{ModeOnce, ModeStream} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			plain := filepath.Join(dir, "plain.bin")
			first := filepath.Join(dir, "first.bin")
			second := filepath.Join(dir, "second.bin")

			writeRandomFile(t, plain, 1<<20)

			priv, pub := newKeyPair(t)
			enc := newEncryptor(t, mode, pub)

			requireNoError(t, enc.Encrypt(plain, first), "first encrypt")
			requireNoError(t, enc.Encrypt(plain, second), "second encrypt")

			if fileHash(t, first) == fileHash(t, second) {
				t.Fatal("encrypting the same plaintext twice produced identical ciphertext")
			}
			if fh1, fh2 := readFileHeader(t, first), readFileHeader(t, second); string(fh1.PublicKey) == string(fh2.PublicKey) {
				t.Fatal("reused encryptor reused its ephemeral public key")
			}

			for _, cipher := range []string{first, second} {
				out := cipher + ".out"
				requireNoError(t, tryDecrypt(priv, cipher, out), "decrypt "+filepath.Base(cipher))
				requireSameContent(t, plain, out)
			}
		})
	}
}

// TestInitMessageFile 验证初始化消息文件的行为。
func TestInitMessageFile(t *testing.T) {
	inWorkDir(t)

	if err := internal.InitMessageFile(); err != nil {
		t.Fatalf("init message file: %v", err)
	}

	got, err := os.ReadFile("message.txt")
	if err != nil {
		t.Fatalf("read message.txt: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("message.txt is empty")
	}

	// 生成的消息文件必须能被后续流程识别并加密
	if err := internal.GenerateKeyPair(true); err != nil {
		t.Fatalf("generate key pair: %v", err)
	}
	if err := internal.Encrypt(); err != nil {
		t.Fatalf("encrypt initialized message file: %v", err)
	}
	if err := internal.Decrypt(); err != nil {
		t.Fatalf("decrypt initialized message file: %v", err)
	}
	requireSameContent(t, "message.txt", "message_decrypted.txt")
}

// TestGenerateKeyPair 验证密钥对生成：默认不覆盖已有文件，显式允许时才覆盖。
func TestGenerateKeyPair(t *testing.T) {
	inWorkDir(t)

	if err := internal.GenerateKeyPair(false); err != nil {
		t.Fatalf("first generate: %v", err)
	}
	privHash, pubHash := fileHash(t, "priv.key"), fileHash(t, "PUB.KEY")

	err := internal.GenerateKeyPair(false)
	if got := errorCode(err); got != utils.ErrFileExist().Code {
		t.Fatalf("second generate without overwrite: error code = %d, want %d", got, utils.ErrFileExist().Code)
	}
	if got := fileHash(t, "priv.key"); got != privHash {
		t.Fatal("private key changed although overwrite was not allowed")
	}

	if err := internal.GenerateKeyPair(true); err != nil {
		t.Fatalf("generate with overwrite: %v", err)
	}
	if privHash == fileHash(t, "priv.key") || pubHash == fileHash(t, "PUB.KEY") {
		t.Fatal("overwrite did not replace the key pair")
	}

	// 新密钥对必须可用
	writeRandomFile(t, "message.txt", 64)
	if err := internal.Encrypt(); err != nil {
		t.Fatalf("encrypt with regenerated key: %v", err)
	}
	if err := internal.Decrypt(); err != nil {
		t.Fatalf("decrypt with regenerated key: %v", err)
	}
	requireSameContent(t, "message.txt", "message_decrypted.txt")
}
