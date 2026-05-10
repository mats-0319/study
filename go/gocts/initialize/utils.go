package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

func backupGenerateFile(fileName string) {
	filePath := token.GeneratorIns.Config.GoDir + fileName

	_, err := os.Stat(filePath)
	if err != nil {
		return
	}

	// file exist, backup old file, backup failed NOT exit
	backupPath := fmt.Sprintf("%s%s/%s-%d.txt", token.GeneratorIns.Config.GoDir,
		utils.GoBackupFolderName, fileName, time.Now().Unix())
	err = os.Rename(filePath, backupPath)
	if err != nil {
		log.Println("backup failed, error:", err)
	}
}
