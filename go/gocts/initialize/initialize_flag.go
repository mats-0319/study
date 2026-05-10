package initialize

import (
	"encoding/json"
	"log"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

// OnInitialize 实现-i参数：写默认配置文件和go接口文件示例文件(demo)
func OnInitialize() {
	writeDefaultConfigFile()
	writeDefaultInitializerFile()
}

func writeDefaultConfigFile() {
	content, err := json.Marshal(token.DefaultGeneratorConfig)
	if err != nil {
		log.Fatalln("json marshal failed, error: ", err)
	}

	fileName := "config.json"
	backupGenerateFile(fileName)
	utils.WriteFile(token.GeneratorIns.Config.GoDir+fileName, content)
}

func writeDefaultInitializerFile() {
	content, err := json.Marshal(DefaultInitializer)
	if err != nil {
		log.Fatalln("json marshal failed, error: ", err)
	}

	fileName := "init.json"
	backupGenerateFile(fileName)
	utils.WriteFile(token.GeneratorIns.Config.GoDir+fileName, content)
}

var DefaultInitializer = &Initializer{
	PackageName: "api",
	Files: []*GoAPIFile{{
		FileName: "demo",
		APIList: []*APIItem{{
			Name: "ListUser",
			URI:  "/user/list",
		}, {
			Name: "CreateUser",
			URI:  "/user/create",
		}},
		EnumList: []*EnumItem{{
			Name:   "UserIdentify",
			Number: 3,
		}},
	}},
}
