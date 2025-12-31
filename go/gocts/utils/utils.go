package utils

import (
	"bytes"
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

// MustSmall make 'str' to small-start camel case, e.g. "Must_small" => "mustSmall"
func MustSmall(str string) string {
	if len(str) < 1 {
		return ""
	}

	str = camelCase(str)

	strBytes := []byte(str)
	if 'A' <= str[0] && str[0] <= 'Z' {
		strBytes[0] = strBytes[0] - 'A' + 'a'
	}

	return string(strBytes)
}

// MustBig make 'str' to big-start camel case, e.g. "must_big" => "MustBig"
func MustBig(str string) string {
	if len(str) < 1 {
		return ""
	}

	str = camelCase(str)

	strBytes := []byte(str)
	if 'a' <= str[0] && str[0] <= 'z' {
		strBytes[0] = strBytes[0] - 'a' + 'A'
	}

	return string(strBytes)
}

func camelCase(str string) string {
	strBytes := make([]byte, 0, len(str))
	underlineFlag := false
	for _, char := range str {
		if char == '_' {
			underlineFlag = true
			continue
		}

		if underlineFlag && 'a' <= char && char <= 'z' {
			strBytes = append(strBytes, byte(char-'a'+'A'))
		} else {
			strBytes = append(strBytes, byte(char))
		}

		underlineFlag = false
	}

	return string(strBytes)
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
		isSep := in(value[index], sep...)

		if !isSep {
			continue
		}

		item := bytes.TrimSpace(value[start:index])
		start = index + 1

		if len(item) > 0 && !bytes.HasPrefix(item, []byte("//")) { // 忽略只有空白字符的结果，忽略注释开头的结果
			res = append(res, item)
		}
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
