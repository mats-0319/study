package test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mats0319/secure_transfer/internal"
	"github.com/mats0319/secure_transfer/internal/lib"
	"github.com/mats0319/secure_transfer/utils"
)

type benchCase struct {
	mode string
	size int64
}

// benchCases 选取两类关键尺寸：
//   - 一次性加密：100B / 4KiB 是「传递一条短消息」的典型场景，1MiB / 8MiB 用于观察随尺寸的变化；
//   - 流式加密：均在一次性加密阈值之上，1MiB 观察帧数较少时的并行开销，
//     8MiB / 32MiB 观察帧数远多于缓冲区数量时的稳态表现（帧缓冲池主要在这里生效）。
//
// 说明：基准把输出文件的写入也算在内，因为这是本工具的真实开销构成。
func benchCases() []benchCase {
	return []benchCase{
		{ModeOnce, 100},
		{ModeOnce, 4 << 10},
		{ModeOnce, 1 << 20},
		{ModeOnce, 8 << 20},
		{ModeStream, 1 << 20},
		{ModeStream, 8 << 20},
		{ModeStream, 32 << 20},
	}
}

// BenchmarkEncrypt 测量加密耗时与分配。每次迭代都会生成新的临时密钥对，
// 与生产路径一致。
func BenchmarkEncrypt(b *testing.B) {
	for _, c := range benchCases() {
		b.Run(fmt.Sprintf("%s/%s", c.mode, humanSize(c.size)), func(b *testing.B) {
			dir := b.TempDir()
			plain := filepath.Join(dir, "plain.bin")
			cipher := filepath.Join(dir, "cipher.bin")

			writeRandomFile(b, plain, c.size)
			_, pub := newKeyPair(b)

			b.SetBytes(c.size)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if e := newEncryptor(b, c.mode, pub).Encrypt(plain, cipher); e != nil {
					b.Fatalf("encrypt: %s", e.String())
				}
			}
		})
	}
}

// BenchmarkDecrypt 测量解密耗时与分配。
// 头部解析与加密只做一次，不计入计时范围。
func BenchmarkDecrypt(b *testing.B) {
	for _, c := range benchCases() {
		b.Run(fmt.Sprintf("%s/%s", c.mode, humanSize(c.size)), func(b *testing.B) {
			dir := b.TempDir()
			plain := filepath.Join(dir, "plain.bin")
			cipher := filepath.Join(dir, "cipher.bin")
			out := filepath.Join(dir, "plain.out")

			writeRandomFile(b, plain, c.size)
			priv, pub := newKeyPair(b)
			encryptFile(b, c.mode, pub, plain, cipher)

			header := readFileHeader(b, cipher)
			cipherSize := fileSize(b, cipher)

			b.SetBytes(c.size)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var dec lib.Decryptor
				if c.mode == ModeOnce {
					dec = lib.NewDecryptorOnce(priv, header, cipherSize)
				} else {
					dec = lib.NewDecryptorStream(priv, header)
				}

				if e := dec.Decrypt(cipher, out); e != nil {
					b.Fatalf("decrypt: %s", e.String())
				}
			}
		})
	}
}

// BenchmarkGenerateKeyPair 测量密钥对生成（X25519 + 编码 + 两份文件写入）。
func BenchmarkGenerateKeyPair(b *testing.B) {
	inWorkDir(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := internal.GenerateKeyPair(true); err != nil {
			b.Fatalf("generate key pair: %v", err)
		}
	}
}

// BenchmarkFirstFile 测量每轮操作都会执行的输入文件匹配（一次目录扫描）。
func BenchmarkFirstFile(b *testing.B) {
	inWorkDir(b)

	writeRandomFile(b, "message.txt", utils.OnceEncMaxSize)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, _, e := utils.FirstFile(utils.OriginFileName); e != nil {
			b.Fatalf("first file: %s", e.String())
		}
	}
}
