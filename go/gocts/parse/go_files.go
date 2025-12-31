package parse

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func TraversalDir() {
	entry, err := os.ReadDir(data.GeneratorIns.Config.GoDir)
	if err != nil {
		log.Fatalln(fmt.Sprintf("read dir %s failed, error: %v", data.GeneratorIns.Config.GoDir, err))
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		var fileInfo fs.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			log.Println("get file info failed, error: ", err)
			continue
		}

		if !strings.HasSuffix(fileInfo.Name(), ".go") {
			continue // ignore not go files
		}

		parseGoFile(fileInfo.Name())
	}

	processStructures()
}

func parseGoFile(filename string) {
	absolutePath := data.GeneratorIns.Config.GoDir + filename
	fileBytes, err := os.ReadFile(absolutePath)
	if err != nil {
		log.Fatalln(fmt.Sprintf("read go file(%s) failed, error: %v", absolutePath, err))
	}

	filename = strings.TrimSuffix(filename, ".go")

	matchRequests(filename, fileBytes)
	matchStructs(filename, fileBytes)
	matchEnums(filename, fileBytes)
}

// processStructures do things:
// 1. set ts type and ts zero value for each field of go 'struct' typ
func processStructures() {
	for _, structureIns := range data.GeneratorIns.Structures {
		if structureIns.Typ.IsEnum {
			continue // ignore enum types, enum has set ts type and ts zero value when matched
		}

		setTsTypeAndZeroValue(structureIns)
	}
}

func setTsTypeAndZeroValue(structureIns *data.StructureItem) {
	for i := 0; i < len(structureIns.Fields); i++ {
		fieldIns := structureIns.Fields[i]
		targetStructure, ok := data.GeneratorIns.Structures[fieldIns.GoType]

		if fieldIns.IsEmbedded {
			if !ok {
				log.Fatalln("unknown embedded struct: ", fieldIns.GoType)
			}
			// 如果一个字段既不具名、也不带有json tag，则将其视作嵌入式结构体，
			// 移除该字段，将该字段所示嵌入式结构体的字段提高到与该字段同级
			structureIns.Fields = append(structureIns.Fields[:i], structureIns.Fields[i+1:]...)
			structureIns.Fields = append(structureIns.Fields, targetStructure.Fields...)
			i--
			continue
		}

		if !ok { // basic type, in type map
			tsBasicType, _ := data.GeneratorIns.TsBasicType[fieldIns.GoType]
			tsZeroValue, _ := data.GeneratorIns.TsZeroValue[tsBasicType]
			fieldIns.TSType = tsBasicType
			fieldIns.TSZeroValue = tsZeroValue
		} else { // self-define type
			fieldIns.TSType = fieldIns.GoType
			// 假设程序扫描到结构体S里的字段F，字段F的类型不是内置类型但具体定义还没扫描到，
			// 也就是说没有办法判断字段F是结构体类型还是枚举类型，这会影响字段F的ts零值设置，
			// 所以我们没有扫描到一个struct就设置其字段的ts类型与零值，而是在全部扫描完成后统一设置
			switch { // if struct field is struct type / enum type
			case targetStructure.Typ.IsStruct:
				fieldIns.TSZeroValue = fmt.Sprintf("new %s()", fieldIns.GoType)
			case targetStructure.Typ.IsEnum:
				fieldIns.TSZeroValue = fmt.Sprintf("%s.%s", targetStructure.Fields[0].TSType, utils.EnumPlaceholder)
			}
		}

		if fieldIns.IsArray {
			fieldIns.TSType = fmt.Sprintf("Array<%s>", fieldIns.TSType)
			fieldIns.TSZeroValue = fmt.Sprintf("new %s()", fieldIns.TSType)
		}
	}
}
