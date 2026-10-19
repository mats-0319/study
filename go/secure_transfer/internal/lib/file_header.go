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
	Salt      []byte // for derive key
	BaseNonce []byte // aes-gcm base nonce
	PublicKey []byte
	AAD       []byte

	Encoded []byte
}

func NewFileHeader(encMethod utils.EncryptMethod, pubKey []byte) (*FileHeader, *utils.Error) {
	fh := &FileHeader{
		EncMethod: encMethod,
		Salt:      utils.GenerateRandomBytes(utils.DeriveKeySaltLength),
		BaseNonce: utils.GenerateRandomBytes(utils.AESBaseNonceLength),
		PublicKey: pubKey,
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

	saltStart := 7
	nonceStart := saltStart + utils.DeriveKeySaltLength
	pubKStart := nonceStart + utils.AESBaseNonceLength
	aadStart := pubKStart + utils.PublicKeyLength
	length := aadStart + utils.FileHeaderAADLength
	if length != utils.FileHeaderLength {
		// 通常情况下不会进入这个分支，写在这里是为了保证后续文件头有修改时，如果整体长度计算规则没有匹配，这里可以第一时间发现
		e = utils.ErrFileHeader().WithParam("mismatched length", length).WithParam("want", utils.FileHeaderLength)
		mlog.Error(e.String())
		return
	}

	fh.Encoded = make([]byte, length)

	copy(fh.Encoded[:2], utils.FileHeaderStart)
	fh.Encoded[2] = fh.EncMethod
	fh.Encoded[3] = byte(utils.DeriveKeySaltLength)
	fh.Encoded[4] = byte(utils.AESBaseNonceLength)
	fh.Encoded[5] = byte(utils.PublicKeyLength)
	fh.Encoded[6] = byte(utils.FileHeaderAADLength)

	copy(fh.Encoded[saltStart:nonceStart], fh.Salt)
	copy(fh.Encoded[nonceStart:pubKStart], fh.BaseNonce)
	copy(fh.Encoded[pubKStart:aadStart], fh.PublicKey)

	fh.AAD = utils.CalcSHA256(fh.Encoded[:aadStart])
	copy(fh.Encoded[aadStart:], fh.AAD)

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

	saltStart := 7
	nonceStart := saltStart + utils.DeriveKeySaltLength
	pubKStart := nonceStart + utils.AESBaseNonceLength
	aadStart := pubKStart + utils.PublicKeyLength

	fh.EncMethod = fileHeaderBytes[2]
	fh.Salt = fileHeaderBytes[saltStart:nonceStart]
	fh.BaseNonce = fileHeaderBytes[nonceStart:pubKStart]
	fh.PublicKey = fileHeaderBytes[pubKStart:aadStart]
	fh.AAD = fileHeaderBytes[aadStart:]

	if !bytes.Equal(fh.AAD, utils.CalcSHA256(fileHeaderBytes[:aadStart])) {
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
	case len(fh.Salt) != utils.DeriveKeySaltLength:
		e = utils.ErrFileHeader().WithParam("salt", fh.Salt)
	case len(fh.BaseNonce) != utils.AESBaseNonceLength:
		e = utils.ErrFileHeader().WithParam("base_nonce", fh.BaseNonce)
	case len(fh.PublicKey) != utils.PublicKeyLength:
		e = utils.ErrFileHeader().WithParam("public_key", fh.PublicKey)
	}

	if e != nil {
		mlog.Error(e.String())
	}

	return
}

func (fh *FileHeader) canDeserialize(data []byte) (e *utils.Error) {
	if len(data) < 7 {
		// 外部已经限制了输入内容长度，通常情况下不会进入该分支，但是因为后续要用到数组索引，所以这里还是写一下
		e = utils.ErrFileHeader().WithParam("fixed length", len(data))
		mlog.Error(e.String())
		return
	}

	switch {
	case string(data[:2]) != utils.FileHeaderStart:
		e = utils.ErrFileHeader().WithParam("start", string(data[:2]))
	case !utils.IsEncryptMethod(data[2]):
		e = utils.ErrFileHeader().WithParam("enc method", data[2])
	case data[3] != utils.DeriveKeySaltLength:
		e = utils.ErrFileHeader().WithParam("salt len", data[3])
	case data[4] != utils.AESBaseNonceLength:
		e = utils.ErrFileHeader().WithParam("base nonce len", data[4])
	case data[5] != utils.PublicKeyLength:
		e = utils.ErrFileHeader().WithParam("pubK len", data[5])
	case data[6] != uint8(utils.FileHeaderAADLength):
		e = utils.ErrFileHeader().WithParam("aad len", data[6])
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
	if fileHeaderLen[0] != utils.FileHeaderLength {
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
