package test

import (
	"os"
	"testing"

	"github.com/mats0319/secure_transfer/utils"
)

// TestFirstFile 覆盖输入文件名匹配规则：`[base].xxx`，且文件名中只能有一个点。
// 这条规则决定了「取哪个文件作为输入」以及「扩展名如何拆分」，
// 因此多一个点、少一个点、大小写、目录等情形都要固定下来。
func TestFirstFile(t *testing.T) {
	cases := []struct {
		name  string
		files []string
		dirs  []string
		size  int64
		want  string // 空表示期望匹配不到
	}{
		{name: "single extension", files: []string{"message.txt"}, size: 10, want: "message.txt"},
		{name: "other extension", files: []string{"message.md"}, size: 3, want: "message.md"},
		{name: "compound extension rejected", files: []string{"message.tar.gz"}, want: ""},
		{name: "dotted stem rejected", files: []string{"my.message.txt"}, want: ""},
		{name: "no extension rejected", files: []string{"message"}, want: ""},
		{name: "case sensitive", files: []string{"Message.txt"}, want: ""},
		{name: "directory ignored", files: []string{"message.txt"}, dirs: []string{"message.d"}, size: 7, want: "message.txt"},
		{name: "first by name wins", files: []string{"message.txt", "message.md"}, size: 1, want: "message.md"},
		{name: "prefix only", files: []string{"message2.txt"}, want: ""},
		{name: "empty dir", want: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inWorkDir(t)

			for _, d := range c.dirs {
				if err := os.Mkdir(d, 0o755); err != nil {
					t.Fatalf("mkdir %s: %v", d, err)
				}
			}
			for _, f := range c.files {
				if err := os.WriteFile(f, make([]byte, c.size), 0o644); err != nil {
					t.Fatalf("write %s: %v", f, err)
				}
			}

			name, size, e := utils.FirstFile(utils.OriginFileName)

			if c.want == "" {
				requireCode(t, e, utils.ErrNoMatchedFile, "FirstFile should not match anything")
				return
			}

			requireNoError(t, e, "FirstFile")
			if name != c.want {
				t.Fatalf("matched %q, want %q", name, c.want)
			}
			if size != c.size {
				t.Fatalf("size = %d, want %d", size, c.size)
			}
		})
	}
}

// TestFirstFileCipherBase 验证解密侧按 `CIPHER` 前缀匹配密文文件。
func TestFirstFileCipherBase(t *testing.T) {
	inWorkDir(t)

	if err := os.WriteFile("CIPHER.TXT", make([]byte, 42), 0o644); err != nil {
		t.Fatalf("write cipher: %v", err)
	}

	name, size, e := utils.FirstFile(utils.EncryptedFileName)
	requireNoError(t, e, "FirstFile(CIPHER)")
	if name != "CIPHER.TXT" || size != 42 {
		t.Fatalf("got (%q, %d), want (\"CIPHER.TXT\", 42)", name, size)
	}
}
