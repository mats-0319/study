package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var emptyErr = errors.New("")

func Log(behavior string, err ...error) {
	if len(err) < 1 {
		log.Println(fmt.Sprintf("> %s success.", behavior))
	} else {
		log.Println(fmt.Sprintf("> %s failed, error: %v", behavior, err[0]))
	}
}

// getFirstFile return full file path matched 'fileName' without extension
func getFirstFile(fileName string) string {
	entry, err := os.ReadDir("./")
	if err != nil {
		Log("Read dir", err)
		return ""
	}

	filePath := "./"
	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		var fileInfo os.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			Log("Read info", err)
			continue
		}

		if !strings.HasPrefix(fileInfo.Name()+".", fileName) {
			continue // ignore files with wrong name
		}

		filePath += fileInfo.Name()
		break
	}

	return filePath
}

func getExtension(filePath string, fileName string) string {
	return strings.TrimPrefix(filePath, "./"+fileName)
}
