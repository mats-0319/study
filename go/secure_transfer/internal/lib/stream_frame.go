package lib

import (
	"fmt"
	"io"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type frame struct {
	index       int32
	isLastFrame bool
	data        []byte
	// buf 是 data 的底层池化缓冲，用于把缓冲归还池中；
	// 为 nil 表示该帧没有占用池化缓冲（例如只承载错误的帧）。
	buf *[]byte
	e   *utils.Error
}

func (f *frame) serialize() {
	if len(f.data) < frameHdrLen+utils.GCMTagLength {
		panic("serialize: frame payload not produced by Seal with 9-byte prefix")
	}

	var fixed [frameHdrLen]byte
	if f.isLastFrame {
		fixed[0] = 1
	} else {
		fixed[0] = 0
	}

	// frame index
	fixed[1] = byte(f.index >> 24)
	fixed[2] = byte(f.index >> 16)
	fixed[3] = byte(f.index >> 8)
	fixed[4] = byte(f.index)

	// payload
	length := int32(len(f.data) - frameHdrLen)
	fixed[5] = byte(length >> 24)
	fixed[6] = byte(length >> 16)
	fixed[7] = byte(length >> 8)
	fixed[8] = byte(length)

	copy(f.data[:frameHdrLen], fixed[:])
}

func (f *frame) deserialize(reader io.Reader, index int) (e *utils.Error) {
	var fixed [frameHdrLen]byte
	n, e := readExactInto(reader, fixed[:])
	if e != nil {
		return
	}
	if n < len(fixed) {
		e = utils.ErrFrame().WithParam("length", n).WithParam("want", len(fixed))
		mlog.Error(e.String())
		return
	}

	switch fixed[0] {
	case 0:
		f.isLastFrame = false
	case 1:
		f.isLastFrame = true
	default:
		e = utils.ErrFrame().WithParam("invalid last flag", fixed[0])
		mlog.Error(e.String())
		return
	}

	f.index = (int32(fixed[1]) << 24) + (int32(fixed[2]) << 16) + (int32(fixed[3]) << 8) + int32(fixed[4])
	if f.index != int32(index) {
		e = utils.ErrFrame().WithParam(fmt.Sprintf("invalid frame index, want: %d, get: %d", index, f.index), "")
		mlog.Error(e.String())
		return
	}

	length := (int32(fixed[5]) << 24) + (int32(fixed[6]) << 16) + (int32(fixed[7]) << 8) + int32(fixed[8])
	if length < utils.GCMTagLength || length > utils.FrameSize+utils.GCMTagLength {
		e = utils.ErrFrame().WithParam("invalid frame payload size", length)
		mlog.Error(e.String())
		return
	}

	// 密文直接读进池化缓冲：大文件下避免「每帧一次约 1MB」的分配
	if f.buf == nil {
		f.buf = getFrameBuf()
	}

	payload := (*f.buf)[:int(length)]
	n, e = readExactInto(reader, payload)
	if e != nil {
		return
	}
	if n != len(payload) {
		e = utils.ErrFrame().WithParam("length", n).WithParam("want", length)
		mlog.Error(e.String())
		return
	}

	f.data = payload

	return
}

func makeNonce(baseNonce []byte, isLastFrame bool, counter int32) []byte {
	nonce := make([]byte, utils.AESBaseNonceLength+1+4)
	copy(nonce[:utils.AESBaseNonceLength], baseNonce)

	if isLastFrame {
		nonce[utils.AESBaseNonceLength] = 1
	} else {
		nonce[utils.AESBaseNonceLength] = 0
	}

	nonce[utils.AESBaseNonceLength+1] = byte(counter >> 24)
	nonce[utils.AESBaseNonceLength+2] = byte(counter >> 16)
	nonce[utils.AESBaseNonceLength+3] = byte(counter >> 8)
	nonce[utils.AESBaseNonceLength+4] = byte(counter)

	return nonce
}
