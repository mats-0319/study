package initialize

import (
	"encoding/json"
	"log"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

// OnInitialize 实现-i参数：写默认配置文件和go接口文件示例文件(demo)
func OnInitialize() {
	writeDefaultConfigFile()
	writeDefaultInitializerFile()
}

func writeDefaultConfigFile() {
	content, err := json.Marshal(data.DefaultGeneratorConfig)
	if err != nil {
		log.Fatalln("json marshal failed, error: ", err)
	}

	filename := "config.json"
	backupGenerateFile(filename)
	utils.WriteFile(data.GeneratorIns.Config.GoDir+filename, content)
}

func writeDefaultInitializerFile() {
	content, err := json.Marshal(DefaultInitializer)
	if err != nil {
		log.Fatalln("json marshal failed, error: ", err)
	}

	filename := "init.json"
	backupGenerateFile(filename)
	utils.WriteFile(data.GeneratorIns.Config.GoDir+filename, content)
}
