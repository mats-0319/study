package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// WriteFile write 'content' into 'file'
func WriteFile(filename string, content []byte) {
	err := os.WriteFile(filename, content, 0777)
	if err != nil {
		log.Fatalln("write file failed, error: ", err)
	}

	log.Println("Generated file: ", filename)
}

// MustSmall make sure first char of 'str' is small-case, e.g. "MustSmall" => "mustSmall"
func MustSmall(str string) string {
	if len(str) < 1 {
		return ""
	}

	bytes := []byte(str)
	if 'A' < str[0] && str[0] < 'Z' {
		bytes[0] = bytes[0] - 'A' + 'a'
	}

	return string(bytes)
}

// MustBig make sure first char of 'str' is big-case, e.g. "mustBig" => "MustBig"
func MustBig(str string) string {
	if len(str) < 1 {
		return ""
	}

	bytes := []byte(str)
	if 'a' < str[0] && str[0] < 'z' {
		bytes[0] = bytes[0] - 'a' + 'A'
	}

	return string(bytes)
}

func MustExistDir(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("'mkdir' on %s failed, error: ", dir), err)
	}
}

// EmptyDir del and re-make dir
func EmptyDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(fmt.Sprintf("rm %s failed, error: ", dir), err)
	}

	MustExistDir(dir)
}

// MustSuffix make sure 'str' is end with 'suffix'
func MustSuffix(str string, suffix string) string {
	if !strings.HasSuffix(str, suffix) {
		str += suffix
	}

	return str
}

// BytesSplit split 'value' into slice with 'sep(s)'
func BytesSplit(value []byte, sep ...byte) [][]byte {
	if len(sep) < 1 {
		return [][]byte{value}
	}

	res := make([][]byte, 0)
	for start, index := 0, 0; index < len(value); index++ {
		if value[index] == '/' && index+1 < len(value) && value[index+1] == '/' { // skip comments ('// xxx')
			res = append(res, value[start:index])
			for ; !in(value[index], '\n'); index++ {
			}
			start = index
		}

		isSep := in(value[index], sep...)

		if !isSep {
			continue
		}

		res = append(res, value[start:index])
		start = index + 1
	}

	return res
}

func in(value byte, list ...byte) bool {
	inFlag := false

	for _, b := range list {
		if b == value {
			inFlag = true
			break
		}
	}

	return inFlag
}
