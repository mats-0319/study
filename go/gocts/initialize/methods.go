package initialize

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func (ins *GoAPIFile) writeFile(packageName string) {
	content := fmt.Sprintf("package %s\n", packageName)

	for i := range ins.EnumList {
		content += ins.EnumList[i].toGo()
	}

	for i := range ins.APIList {
		content += ins.APIList[i].toGo()
	}

	ins.FileName = utils.MustSuffix(ins.FileName, ".go")

	filePath := data.GeneratorIns.Config.GoDir + ins.FileName

	_, err := os.Stat(filePath)
	if err == nil { // file exist, backup old file, for dev
		backupPath := data.GeneratorIns.Config.GoDir + fmt.Sprintf("backup/%s-%d.txt",
			strings.TrimSuffix(ins.FileName, ".go"), time.Now().Unix())
		if err = os.Rename(filePath, backupPath); err != nil {
			log.Fatalln("backup failed, error:", err)
		}
	}

	utils.WriteFile(filePath, []byte(content))
}

func (ins *APIItem) toGo() string {
	res := `
const URI_{{ $apiName }} = "{{ $apiURI }}"

type {{ $apiName }}{{ $reqSuffix }} struct {}

type {{ $apiName }}{{ $resSuffix }} struct {}
`

	res = strings.ReplaceAll(res, "{{ $apiName }}", ins.Name)
	res = strings.ReplaceAll(res, "{{ $apiURI }}", ins.URI)
	res = strings.ReplaceAll(res, "{{ $reqSuffix }}", data.GeneratorIns.Config.RequestStructureSuffix)
	res = strings.ReplaceAll(res, "{{ $resSuffix }}", data.GeneratorIns.Config.ResponseStructureSuffix)

	return res
}

func (ins *EnumItem) toGo() string {
	if ins.Number < 3 { // enum type usually larger than 2
		ins.Number = 3
	} else if ins.Number > 1<<7-1 {
		ins.Number = 1<<7 - 1
	}

	enumUnitsStr := ""
	for i := range ins.Number {
		enumUnitsStr += "\n{{ $indentation }}{{ $enumName }}_Value" + strconv.Itoa(i) + " {{ $enumName }} = " + strconv.Itoa(i)
	}
	enumUnitsStr += "\n"

	enumItemStr := `
type {{ $enumName }} int8

const ({{ $enumUnits }})
`

	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $enumUnits }}", enumUnitsStr)
	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $enumName }}", ins.Name)
	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return enumItemStr
}
