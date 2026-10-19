package lib

import (
	"errors"
	"io"
	"sync"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

// frameHdrLen 是单帧帧头长度：末帧标识(1) + 帧索引(4) + 有效内容长度(4)
const frameHdrLen = 9

// frameBufSize 是池化帧缓冲的容量。
// 一帧最大为：帧头 + 一个完整明文帧 + GCM tag，
// 因为明文 payload 被限制在 utils.FrameSize 以内，密文最多再多一个 tag。
const frameBufSize = frameHdrLen + utils.FrameSize + utils.GCMTagLength

// frameBufPool 复用帧缓冲，避免大文件流式处理时「每帧分配约 1MB」。
//
// 这里存 *[]byte 而不是 []byte 是有意为之：把切片放进 sync.Pool 会把切片头
// 装箱进一个 interface，每次 Put 都会产生一次堆分配，反而抵消了池化的收益；
// 放指针则不会分配。
var frameBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, frameBufSize)
		return &b
	},
}

// getFrameBuf 从池中取一块帧缓冲。
func getFrameBuf() *[]byte {
	b := frameBufPool.Get().(*[]byte)
	if cap(*b) < frameBufSize {
		// 防御分支：正常情况下池中只会是 frameBufSize 的缓冲
		nb := make([]byte, frameBufSize)
		return &nb
	}

	return b
}

// putFrameBuf 把缓冲归还池中。
// 调用后 b 以及由它派生的任何切片都不再有效。
func putFrameBuf(b *[]byte) {
	if b == nil {
		return
	}
	if cap(*b) < frameBufSize {
		return // 非本池的缓冲，直接丢弃交给 GC
	}

	*b = (*b)[:frameBufSize]
	frameBufPool.Put(b)
}

// takeBuf 把 b 设为 frame 的底层缓冲，并归还 frame 原来持有的缓冲。
// 调用前应先把 f.data 指向 b 内部的视图。
func (f *frame) takeBuf(b *[]byte) {
	if f.buf != nil && f.buf != b {
		putFrameBuf(f.buf)
	}

	f.buf = b
}

// releaseBuf 归还 frame 持有的缓冲。
// 调用后 f.data（以及由它派生的切片）都不再有效。
func (f *frame) releaseBuf() {
	if f.buf != nil {
		putFrameBuf(f.buf)
		f.buf = nil
	}

	f.data = nil
}

// newPooledFrame 返回一个已挂好池化缓冲的 frame。
func newPooledFrame() *frame {
	return &frame{buf: getFrameBuf()}
}

// readExactInto 把 reader 的内容读满 buf，返回实际读到的字节数。
// 读到 EOF 不算错误：调用方用返回的 n 与 len(buf) 比较来判断是否为最后一帧。
func readExactInto(reader io.Reader, buf []byte) (n int, e *utils.Error) {
	for n < len(buf) { // io.Reader.Read 一次可能只读一部分，循环保证读满或读到结尾
		read, err := reader.Read(buf[n:])
		n += read

		if err != nil {
			if !errors.Is(err, io.EOF) { // 读到结尾不视为错误，正常返回已读内容
				e = utils.ErrReadExact().WithCause(err)
				mlog.Error(e.String())
			}

			return
		}
	}

	return
}

// readExact 分配 length 字节并读满，返回实际读到的内容（可能不足 length）。
func readExact(reader io.Reader, length int) (res []byte, e *utils.Error) {
	buf := make([]byte, length)

	n, e := readExactInto(reader, buf)

	return buf[:n], e
}
