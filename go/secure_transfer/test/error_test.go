package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mats0319/secure_transfer/internal"
	"github.com/mats0319/secure_transfer/utils"
)

// TestDecryptTamperedCipher 覆盖密文被篡改 / 截断 / 头部被破坏时的行为：
// 必须返回错误、不能 panic（历史上负的帧长度会触发不可恢复的 panic）。
func TestDecryptTamperedCipher(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.bin")
	cipher := filepath.Join(dir, "cipher.bin")

	const plainSize = 5 * utils.FrameSize // 5 个满帧 + 1 个空末帧
	writeRandomFile(t, plain, plainSize)

	priv, pub := newKeyPair(t)
	encryptFile(t, ModeStream, pub, plain, cipher)

	raw, err := os.ReadFile(cipher)
	if err != nil {
		t.Fatalf("read cipher: %v", err)
	}

	const fullFrame = frameHeaderLen + utils.FrameSize + utils.GCMTagLength
	frameStart := func(i int) int { return cipherHeaderLen + i*fullFrame }

	// 文件头内部布局：'ST'(2) + enc method(1) + 四个长度字段(4) + salt + base nonce + pubK + aad
	const aadStart = 2 + 1 + 4 + utils.DeriveKeySaltLength + utils.AESBaseNonceLength + utils.PublicKeyLength

	clone := func(mutate func(b []byte)) []byte {
		b := append([]byte(nil), raw...)
		mutate(b)
		return b
	}

	cases := []struct {
		name string
		data []byte
	}{
		{"truncate at frame boundary", raw[:frameStart(3)]},
		{"truncate inside frame payload", raw[:frameStart(2)+500]},
		{"truncate inside file header", raw[:40]},
		{"empty file", nil},
		{"flip ciphertext byte", clone(func(b []byte) { b[frameStart(2)+frameHeaderLen+10] ^= 0xFF })},
		{"flip frame header flag", clone(func(b []byte) { b[frameStart(1)] = 0xFF })},
		{"rewrite frame index", clone(func(b []byte) { copy(b[frameStart(2)+1:], []byte{0, 0, 0, 7}) })},
		{"mark early frame as last", clone(func(b []byte) { b[frameStart(1)] = 1 })},
		{"frame length zero", clone(func(b []byte) { copy(b[frameStart(2)+5:], []byte{0, 0, 0, 0}) })},
		{"frame length negative", clone(func(b []byte) { copy(b[frameStart(2)+5:], []byte{0xFF, 0xFF, 0xFF, 0xFF}) })},
		{"frame length too large", clone(func(b []byte) { copy(b[frameStart(2)+5:], []byte{0x7F, 0xFF, 0xFF, 0xFF}) })},
		{"broken magic", clone(func(b []byte) { copy(b[1:], []byte("XX")) })},
		{"unknown enc method", clone(func(b []byte) { b[3] = 9 })},
		{"wrong header length byte", clone(func(b []byte) { b[0] = 42 })},
		{"flip aad hash", clone(func(b []byte) { b[1+aadStart] ^= 0xFF })},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			badDir := t.TempDir()
			bad := filepath.Join(badDir, "tampered.bin")
			out := filepath.Join(badDir, "out.bin")

			if err := os.WriteFile(bad, c.data, 0o644); err != nil {
				t.Fatalf("write tampered cipher: %v", err)
			}

			var decErr *utils.Error
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("decrypt panicked: %v", r)
					}
				}()
				decErr = tryDecrypt(priv, bad, out)
			}()

			if decErr == nil {
				t.Fatal("tampered cipher decrypted successfully")
			}
		})
	}
}

// TestDecryptIntactStreamCipher 确认同一份未篡改密文可以正常解密，
// 保证上面的篡改用例失败确实由篡改引起。
func TestDecryptIntactStreamCipher(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.bin")
	cipher := filepath.Join(dir, "cipher.bin")
	out := filepath.Join(dir, "plain.out")

	writeRandomFile(t, plain, 3*utils.FrameSize)

	priv, pub := newKeyPair(t)
	encryptFile(t, ModeStream, pub, plain, cipher)

	requireNoError(t, tryDecrypt(priv, cipher, out), "decrypt intact cipher")
	requireSameContent(t, plain, out)
}

// TestDecryptWrongPrivateKey 验证用非目标接收方的私钥无法解密。
func TestDecryptWrongPrivateKey(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.bin")
	cipher := filepath.Join(dir, "cipher.bin")
	out := filepath.Join(dir, "plain.out")

	writeRandomFile(t, plain, 1024)

	_, pub := newKeyPair(t)
	encryptFile(t, ModeOnce, pub, plain, cipher)

	wrongPriv, _ := newKeyPair(t)
	requireCode(t, tryDecrypt(wrongPriv, cipher, out), utils.ErrAESDecrypt, "decrypt with wrong private key")
}

// TestOnceCipherOverLimit 验证一次性密文的解密端内存上限仍然生效：
// 伪造一个超大的「一次性加密」密文时应当在读取正文之前就被拒绝，
// 而不是按文件大小分配内存。
func TestOnceCipherOverLimit(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.bin")
	cipher := filepath.Join(dir, "cipher.bin")
	out := filepath.Join(dir, "plain.out")

	writeRandomFile(t, plain, 1024)

	priv, pub := newKeyPair(t)
	encryptFile(t, ModeOnce, pub, plain, cipher)

	// 补零到「上限 + 1」，让密文长度刚好越过解密端允许的最大值
	limit := int64(cipherHeaderLen+utils.OnceEncMaxSize+utils.GCMTagLength) + 1
	padLen := limit - fileSize(t, cipher)
	if padLen <= 0 {
		t.Fatalf("cipher already larger than the limit: %d", fileSize(t, cipher))
	}

	f, err := os.OpenFile(cipher, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("open cipher for padding: %v", err)
	}
	if _, err := f.Write(make([]byte, padLen)); err != nil {
		t.Fatalf("pad cipher: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close padded cipher: %v", err)
	}

	requireCode(t, tryDecrypt(priv, cipher, out), utils.ErrEncryptedFile, "decrypt oversized once cipher")
}

// TestDecryptRemovesPartialOutput 验证生产入口在解密失败时会删除半成品文件。
// 破坏末帧的 tag，使前面的帧已经解密并写盘，随后才失败。
func TestDecryptRemovesPartialOutput(t *testing.T) {
	inWorkDir(t)

	if err := internal.GenerateKeyPair(true); err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	// 必须大于一次性加密阈值，生产路径才会选择流式加密
	writeRandomFile(t, "message.txt", utils.OnceEncMaxSize+1)

	if err := internal.Encrypt(); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	raw, err := os.ReadFile("CIPHER.TXT")
	if err != nil {
		t.Fatalf("read cipher: %v", err)
	}
	raw[len(raw)-1] ^= 0xFF // 破坏末帧 tag
	if err := os.WriteFile("CIPHER.TXT", raw, 0o644); err != nil {
		t.Fatalf("write corrupted cipher: %v", err)
	}

	if err := internal.Decrypt(); err == nil {
		t.Fatal("decrypting a corrupted cipher succeeded")
	}

	if _, err := os.Stat("message_decrypted.txt"); err == nil {
		t.Fatal("partial plaintext file was left behind after a failed decrypt")
	}
}

// TestMissingInputFile 验证工作目录里没有输入文件时返回 ErrNoMatchedFile。
func TestMissingInputFile(t *testing.T) {
	inWorkDir(t)

	if err := internal.GenerateKeyPair(true); err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	if got := errorCode(internal.Encrypt()); got != utils.ErrNoMatchedFile().Code {
		t.Fatalf("encrypt without message file: error code = %d, want %d", got, utils.ErrNoMatchedFile().Code)
	}
	if got := errorCode(internal.Decrypt()); got != utils.ErrNoMatchedFile().Code {
		t.Fatalf("decrypt without cipher file: error code = %d, want %d", got, utils.ErrNoMatchedFile().Code)
	}
}
