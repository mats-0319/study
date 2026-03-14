package initialize

import (
	"fmt"
	"strings"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

type Initializer struct {
	PackageName string       `json:"package_name"` // 因为同文件夹下.go文件的包名应统一，这里就提出来了
	Files       []*GoAPIFile `json:"files"`
}

var InitializerIns = &Initializer{}

type GoAPIFile struct {
	FileName string      `json:"file_name"`
	APIList  []*APIItem  `json:"api_list"`
	EnumList []*EnumItem `json:"enum_list"`
}

func (ins *GoAPIFile) writeFile(packageName string) {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("package %s\n", packageName))

	for i := range ins.EnumList {
		content.WriteString(ins.EnumList[i].toGo())
	}

	for i := range ins.APIList {
		content.WriteString(ins.APIList[i].toGo())
	}

	ins.FileName = utils.MustSuffix(ins.FileName, ".go")

	backupGenerateFile(ins.FileName)

	filePath := token.GeneratorIns.Config.GoDir + ins.FileName
	utils.WriteFile(filePath, content.String())
}

type APIItem struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func (ins *APIItem) toGo() string {
	res := `
const URI_{{ $apiName }} = "{{ $apiURI }}"

type {{ $apiName }}{{ $reqSuffix }} struct {}

type {{ $apiName }}{{ $resSuffix }} struct {}
`

	res = strings.ReplaceAll(res, "{{ $apiName }}", ins.Name)
	res = strings.ReplaceAll(res, "{{ $apiURI }}", ins.URI)
	res = strings.ReplaceAll(res, "{{ $reqSuffix }}", token.GeneratorIns.Config.RequestStructureSuffix)
	res = strings.ReplaceAll(res, "{{ $resSuffix }}", token.GeneratorIns.Config.ResponseStructureSuffix)

	return res
}

type EnumItem struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func (ins *EnumItem) toGo() string {
	ins.Number = min(ins.Number, 20) // not too much

	enumUnitsStr := ""
	if ins.Number > 0 {
		enumUnitsStr = "\n"
	}
	for i := range ins.Number {
		enumUnitsStr += fmt.Sprintf("{{ $indentation }}{{ $enumName }}_Value%d {{ $enumName }} = %d\n", i+1, i+1)
	}

	enumItemStr := `
type {{ $enumName }} int8

const ({{ $enumUnits }})
`

	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $enumUnits }}", enumUnitsStr)
	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $enumName }}", ins.Name)
	enumItemStr = strings.ReplaceAll(enumItemStr, "{{ $indentation }}", token.GeneratorIns.IndentationStr)

	return enumItemStr
}
