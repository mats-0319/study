package lib

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type FileHeader struct {
	EncMethod utils.EncryptMethod
	PublicKey []byte
	BaseNonce []byte // aes-gcm base nonce
	Salt      []byte // for derive key
	AAD       []byte

	Encoded []byte
}

func NewFileHeader(encMethod utils.EncryptMethod, pubKey []byte) (*FileHeader, *utils.Error) {
	fh := &FileHeader{
		EncMethod: encMethod,
		PublicKey: pubKey,
		BaseNonce: utils.GenerateRandomBytes(utils.AESBaseNonceLength),
		Salt:      utils.GenerateRandomBytes(utils.DeriveKeySaltLength),
	}

	e := fh.serialize()
	if e != nil {
		return nil, e
	}

	return fh, nil
}

func (fh *FileHeader) serialize() (e *utils.Error) {
	e = fh.canSerialize()
	if e != nil {
		return
	}

	fh.Encoded = make([]byte, utils.FileHeaderLength-1)

	copy(fh.Encoded[:2], utils.FileHeaderFlag)
	fh.Encoded[2] = fh.EncMethod

	start := 3
	fh.Encoded[start] = utils.PublicKeyLength
	copy(fh.Encoded[start+1:start+1+utils.PublicKeyLength], fh.PublicKey)

	start += 1 + utils.PublicKeyLength
	fh.Encoded[start] = utils.AESBaseNonceLength
	copy(fh.Encoded[start+1:start+1+utils.AESBaseNonceLength], fh.BaseNonce)

	start += 1 + utils.AESBaseNonceLength
	fh.Encoded[start] = utils.DeriveKeySaltLength
	copy(fh.Encoded[start+1:start+1+utils.DeriveKeySaltLength], fh.Salt)

	start += 1 + utils.DeriveKeySaltLength
	fh.AAD = utils.CalcSHA256(fh.Encoded[:start])
	fh.Encoded[start] = utils.FileHeaderAADLength
	copy(fh.Encoded[start+1:], fh.AAD)

	return
}

func (fh *FileHeader) Deserialize(encFile string) (e *utils.Error) {
	file, err := os.Open(encFile)
	if err != nil {
		e = utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, utils.FileHeaderLength)

	fileHeaderBytes, e := readFileHeader(reader)
	if e != nil {
		return
	}

	// deserialize
	e = fh.canDeserialize(fileHeaderBytes)
	if e != nil {
		return
	}

	fh.EncMethod = fileHeaderBytes[2]

	index := 3
	length := int(fileHeaderBytes[index])
	fh.PublicKey = fileHeaderBytes[index+1 : index+1+length]

	index += 1 + length
	length = int(fileHeaderBytes[index])
	fh.BaseNonce = fileHeaderBytes[index+1 : index+1+length]

	index += 1 + length
	length = int(fileHeaderBytes[index])
	fh.Salt = fileHeaderBytes[index+1 : index+1+length]

	index += 1 + length
	length = int(fileHeaderBytes[index])
	fh.AAD = fileHeaderBytes[index+1 : index+1+length]

	if !bytes.Equal(fh.AAD, utils.CalcSHA256(fileHeaderBytes[:index])) {
		e = utils.ErrFileHeader().WithParam("hash not matched", "")
		mlog.Error(e.String())
		return
	}

	return
}

func (fh *FileHeader) canSerialize() (e *utils.Error) {
	switch {
	case !utils.IsEncryptMethod(fh.EncMethod):
		e = utils.ErrFileHeader().WithParam("enc method", fh.EncMethod)
	case len(fh.PublicKey) != utils.PublicKeyLength:
		e = utils.ErrFileHeader().WithParam("public_key", fh.PublicKey)
	case len(fh.BaseNonce) != utils.AESBaseNonceLength:
		e = utils.ErrFileHeader().WithParam("base_nonce", fh.BaseNonce)
	case len(fh.Salt) != utils.DeriveKeySaltLength:
		e = utils.ErrFileHeader().WithParam("salt", fh.Salt)
	}

	if e != nil {
		mlog.Error(e.String())
	}

	return
}

func (fh *FileHeader) canDeserialize(data []byte) (e *utils.Error) {
	switch {
	case len(data) != utils.FileHeaderLength-1:
		e = utils.ErrFileHeader().WithParam("fixed length", len(data))
	case string(data[:2]) != utils.FileHeaderFlag:
		e = utils.ErrFileHeader().WithParam("flag", string(data[:2]))
	case !utils.IsEncryptMethod(data[2]):
		e = utils.ErrFileHeader().WithParam("enc method", data[2])
	}

	if e != nil {
		mlog.Error(e.String())
	}

	return
}

func readFileHeader(reader io.Reader) (fileHeaderBytes []byte, e *utils.Error) {
	fileHeaderLen, e := readExact(reader, 1)
	if e != nil {
		return
	}
	if len(fileHeaderLen) < 1 {
		// 空文件：readExact 读到 EOF 且不视为错误，这里必须显式拒绝，
		// 否则下面索引 fileHeaderLen[0] 会越界 panic
		e = utils.ErrFileHeader().WithParam("length flag", "empty file")
		mlog.Error(e.String())
		return
	}
	if fileHeaderLen[0] != utils.FileHeaderLength-1 {
		e = utils.ErrFileHeader().
			WithParam("length flag", fileHeaderLen[0]).
			WithParam("want", utils.FileHeaderLength)
		mlog.Error(e.String())
		return
	}

	fileHeaderBytes, e = readExact(reader, int(fileHeaderLen[0]))
	if e != nil {
		return
	}
	if len(fileHeaderBytes) != int(fileHeaderLen[0]) {
		e = utils.ErrFileHeader().
			WithParam("length", len(fileHeaderBytes)).
			WithParam("want", fileHeaderLen[0])
		mlog.Error(e.String())
		return
	}

	return
}
