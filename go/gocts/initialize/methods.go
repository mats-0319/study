package initialize

import (
	"fmt"
	"strings"

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

	backupGenerateFile(ins.FileName)

	filePath := data.GeneratorIns.Config.GoDir + ins.FileName
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
	if ins.Number > 1<<7-1 {
		ins.Number = 1<<7 - 1
	}

	enumUnitsStr := ""
	if ins.Number > 0 {
		enumUnitsStr = "\n"
	}
	for i := range ins.Number { // 'ins.Number' < 0 is ok
		enumUnitsStr += fmt.Sprintf("{{ $indentation }}{{ $enumName }}_Value%d {{ $enumName }} = %d\n", i, i)
	}

	enumItemStr := `
type {{ $enumName }} int8

const ({{ $enumUnits }})
`

	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $enumUnits }}", enumUnitsStr)
	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $enumName }}", ins.Name)
	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return enumItemStr
}
