package test

// 本包的测试约定：
//
//   - internal.Encrypt / internal.Decrypt / internal.GenerateKeyPair / internal.InitMessageFile
//     都使用「当前工作目录 + 固定文件名」，因此需要切换工作目录，用 inWorkDir 完成；
//   - os.Chdir 是进程级状态，所以本包任何测试都不能调用 t.Parallel()；
//   - 测试数据一律在临时目录生成，不依赖 dd 等外部命令（可在 Windows / Linux / macOS 上运行）。

import (
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/mats0319/secure_transfer/internal/lib"
	"github.com/mats0319/secure_transfer/utils"
)

const (
	// ModeOnce / ModeStream 是测试里表示两种加密方式的字符串标识。
	ModeOnce   = "once"
	ModeStream = "stream"

	// frameHeaderLen 是流式密文单帧帧头的长度：末帧标识(1) + 帧索引(4) + 有效内容长度(4)。
	// 该值属于加密文件编码规则，故意在测试中独立写死，以便格式一旦变动测试立刻失败。
	frameHeaderLen = 9

	// cipherHeaderLen 是密文文件头长度：长度前缀(1) + 文件头(94)。
	cipherHeaderLen = 1 + 94
)

// inWorkDir 切换到新的临时工作目录，并在测试结束后切回。
// 切回动作通过 Cleanup 注册在 t.TempDir 的删除动作之后，因此会先于删除执行
// （Windows 下不能删除仍为当前工作目录的目录）。
func inWorkDir(tb testing.TB) string {
	tb.Helper()

	dir := tb.TempDir()

	old, err := os.Getwd()
	if err != nil {
		tb.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		tb.Fatalf("chdir %s: %v", dir, err)
	}

	tb.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			tb.Errorf("restore working directory: %v", err)
		}
	})

	return dir
}

// writeRandomFile 在 Go 内生成指定大小的伪随机文件。
func writeRandomFile(tb testing.TB, path string, size int64) {
	tb.Helper()

	f, err := os.Create(path)
	if err != nil {
		tb.Fatalf("create %s: %v", path, err)
	}
	defer func() { _ = f.Close() }() // 出错时兜底关闭

	src := rand.NewChaCha8([32]byte{1, 2, 3}) // 固定种子，保证可复现
	buf := make([]byte, 1<<20)
	for written := int64(0); written < size; {
		n := int64(len(buf))
		if size-written < n {
			n = size - written
		}

		chunk := buf[:n]
		for i := int64(0); i < n; i += 8 {
			v := src.Uint64()
			for j := int64(0); j < 8 && i+j < n; j++ {
				chunk[i+j] = byte(v >> (8 * j))
			}
		}

		if _, err := f.Write(chunk); err != nil {
			tb.Fatalf("write %s: %v", path, err)
		}
		written += n
	}

	if err := f.Close(); err != nil {
		tb.Fatalf("close %s: %v", path, err)
	}
}

func fileSize(tb testing.TB, path string) int64 {
	tb.Helper()

	fi, err := os.Stat(path)
	if err != nil {
		tb.Fatalf("stat %s: %v", path, err)
	}

	return fi.Size()
}

func fileHash(tb testing.TB, path string) string {
	tb.Helper()

	f, err := os.Open(path)
	if err != nil {
		tb.Fatalf("open %s: %v", path, err)
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		tb.Fatalf("hash %s: %v", path, err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// requireSameContent 断言两个文件内容完全一致。
func requireSameContent(tb testing.TB, want, got string) {
	tb.Helper()

	if wantSize, gotSize := fileSize(tb, want), fileSize(tb, got); wantSize != gotSize {
		tb.Fatalf("size mismatch: %s = %d bytes, %s = %d bytes", want, wantSize, got, gotSize)
	}
	if wantHash, gotHash := fileHash(tb, want), fileHash(tb, got); wantHash != gotHash {
		tb.Fatalf("content mismatch: %s and %s have different sha256", want, got)
	}
}

func newKeyPair(tb testing.TB) (*ecdh.PrivateKey, *ecdh.PublicKey) {
	tb.Helper()

	priv, err := utils.Curve.GenerateKey(nil)
	if err != nil {
		tb.Fatalf("generate key pair: %v", err)
	}

	return priv, priv.PublicKey()
}

func newEncryptor(tb testing.TB, mode string, pub *ecdh.PublicKey) lib.Encryptor {
	tb.Helper()

	switch mode {
	case ModeOnce:
		return lib.NewEncryptorOnce(pub)
	case ModeStream:
		return lib.NewEncryptorStream(pub)
	default:
		tb.Fatalf("unknown encrypt mode %q", mode)
		return nil
	}
}

func encryptFile(tb testing.TB, mode string, pub *ecdh.PublicKey, plain, cipher string) {
	tb.Helper()

	if e := newEncryptor(tb, mode, pub).Encrypt(plain, cipher); e != nil {
		tb.Fatalf("encrypt %s -> %s: %s", filepath.Base(plain), filepath.Base(cipher), e.String())
	}
}

// tryDecrypt 按密文头记录的加密方式选择解密器，返回底层错误而不判定测试失败。
func tryDecrypt(priv *ecdh.PrivateKey, cipher, out string) *utils.Error {
	fh := &lib.FileHeader{}
	if e := fh.Deserialize(cipher); e != nil {
		return e
	}

	fi, err := os.Stat(cipher)
	if err != nil {
		return utils.ErrOpenEncFile().WithCause(err)
	}

	switch fh.EncMethod {
	case utils.EncryptMethod_Once:
		return lib.NewDecryptorOnce(priv, fh, fi.Size()).Decrypt(cipher, out)
	case utils.EncryptMethod_Stream:
		return lib.NewDecryptorStream(priv, fh).Decrypt(cipher, out)
	default:
		return utils.ErrFileHeader().WithParam("enc method", fh.EncMethod)
	}
}

func readFileHeader(tb testing.TB, cipher string) *lib.FileHeader {
	tb.Helper()

	fh := &lib.FileHeader{}
	if e := fh.Deserialize(cipher); e != nil {
		tb.Fatalf("deserialize header of %s: %s", filepath.Base(cipher), e.String())
	}

	return fh
}

// modeOf 返回密文头里记录的加密方式标识。
func modeOf(tb testing.TB, cipher string) string {
	tb.Helper()

	switch readFileHeader(tb, cipher).EncMethod {
	case utils.EncryptMethod_Once:
		return ModeOnce
	case utils.EncryptMethod_Stream:
		return ModeStream
	default:
		tb.Fatalf("unknown enc method byte in header")
		return ""
	}
}

// requireCipherSize 校验密文长度符合编码规则。
func requireCipherSize(tb testing.TB, mode, cipher string, plainSize int64) {
	tb.Helper()

	var want int64
	switch mode {
	case ModeOnce:
		want = cipherHeaderLen + plainSize + utils.GCMTagLength
	case ModeStream:
		// 帧数 = 整除商 + 1：恰好整除时末尾会补一个空帧，不整除时最后一帧不足一帧且带末帧标识
		frames := plainSize/utils.FrameSize + 1
		want = cipherHeaderLen + frames*(frameHeaderLen+utils.GCMTagLength) + plainSize
	}

	if got := fileSize(tb, cipher); got != want {
		tb.Fatalf("%s cipher size = %d, want %d (plain = %d)", mode, got, want, plainSize)
	}
}

func requireNoError(tb testing.TB, e *utils.Error, what string) {
	tb.Helper()

	if e != nil {
		tb.Fatalf("%s: unexpected error: %s", what, e.String())
	}
}

// requireCode 断言错误码与 want() 生成的错误一致。
func requireCode(tb testing.TB, e *utils.Error, want func() *utils.Error, what string) {
	tb.Helper()

	if e == nil {
		tb.Fatalf("%s: expected error %d, got success", what, want().Code)
	}
	if e.Code != want().Code {
		tb.Fatalf("%s: expected error code %d, got %d (%s)", what, want().Code, e.Code, e.Message)
	}
}

// errorCode 从 error 中取出 utils.Error 的错误码，非 utils.Error 返回 0。
func errorCode(err error) int {
	var e *utils.Error
	if errors.As(err, &e) && e != nil {
		return e.Code
	}

	return 0
}

// humanSize 生成可读的尺寸标识，用于子测试 / 基准测试命名。
func humanSize(n int64) string {
	switch {
	case n == 0:
		return "0B"
	case n%(1<<20) == 0:
		return fmt.Sprintf("%dMiB", n>>20)
	case n < 1<<20 && n%(1<<10) == 0:
		return fmt.Sprintf("%dKiB", n>>10)
	default:
		return fmt.Sprintf("%dB", n)
	}
}
