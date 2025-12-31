package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func ExecLog(behavior string, err ...error) {
	if len(err) < 1 {
		fmt.Println(fmt.Sprintf("> %s success.", behavior))
	} else {
		log.Println(fmt.Sprintf("> %s failed, error: %v", behavior, err[0]))
	}
}

// getFirstFile return full file path matched 'fileName' without extension
func getFirstFile(fileName string) (string, error) {
	entry, err := os.ReadDir("./")
	if err != nil {
		ExecLog("Read dir", err)
		return "", err
	}

	filePath := "./"
	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		var fileInfo os.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			ExecLog("Read info", err)
			continue
		}

		if !strings.HasPrefix(fileInfo.Name()+".", fileName) {
			continue // ignore files with wrong name
		}

		filePath += fileInfo.Name()
		break
	}

	return filePath, nil
}

func getExtension(filePath string, fileName string) string {
	return strings.TrimPrefix(filePath, "./"+fileName)
}
