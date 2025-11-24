package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func backupGenerateFile(filename string) {
	filePath := data.GeneratorIns.Config.GoDir + filename

	_, err := os.Stat(filePath)
	if err != nil {
		return
	}

	// file exist, backup old file, backup failed NOT exit
	backupPath := fmt.Sprintf("%s%s/%s-%d.txt", data.GeneratorIns.Config.GoDir,
		utils.GoBackupFolderName, filename, time.Now().Unix())
	err = os.Rename(filePath, backupPath)
	if err != nil {
		log.Println("backup failed, error:", err)
	}
}
