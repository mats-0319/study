package utils

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"os"
	"strings"

	mlog "github.com/mats0319/secure_transfer/utils/log"
)

var Curve = ecdh.X25519()

// FirstFile get first file which name matched 'file base name'
// - fileBaseName: file name without extension
func FirstFile(fileBaseName string) (fileName string, fileSize int64, e *Error) {
	entry, err := os.ReadDir("./")
	if err != nil {
		e = ErrReadDir().WithCause(err)
		mlog.Error(e.String())
		return
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		// match file name, require: '[fileBaseName].xxx',
		// 强制限制文件名格式中只能有一个点，是为了方便取扩展名、往返名称对称
		entryName := entry[i].Name()
		if !strings.HasPrefix(entryName, fileBaseName+".") ||
			!(strings.Index(entryName, ".") == strings.LastIndex(entryName, ".")) {
			continue
		}

		fileInfo, err := entry[i].Info()
		if err != nil {
			e = ErrGetFileInfo().WithCause(err).WithParam("name", entry[i].Name())
			mlog.Error(e.String())
			continue
		}

		fileName = entryName
		fileSize = fileInfo.Size()

		break
	}

	if fileName == "" {
		e = ErrNoMatchedFile()
	}

	return
}

func GenerateRandomBytes(length int) []byte {
	data := make([]byte, length)
	_, _ = rand.Read(data) // err always nil

	return data
}

func CalcSHA256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)

	return hasher.Sum(nil)
}
