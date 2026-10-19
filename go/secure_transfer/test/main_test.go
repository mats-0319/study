package test

import (
	"log/slog"
	"os"
	"testing"
)

// TestMain 把 slog 的默认 logger 换成丢弃处理器。
//
// 库层在出错时会用 mlog 打印完整调用栈，而本包的用例会大量故意触发错误路径
// （篡改、截断、错误密钥等），保留默认输出会让测试日志几乎不可读。
// 真正的错误信息由各用例的 Fatalf 呈现（其中包含 utils.Error.String()）。
func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.DiscardHandler))

	os.Exit(m.Run())
}
