package parse

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

var (
	RequestRE     = regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*["']([\w/-]+)["']`)
	StructRE      = regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	StructFieldRE = regexp.MustCompile(`\s*(?:\w+\s+)?([\[\]*\w]+)(?:\s+.*?json:["'](\w+)["'].*?)?`)
	EnumRE        = regexp.MustCompile(`type\s+(\w+)\s*=?\s*\w+\s+const\s*\(([^)]*)\)`)
	EnumUnitRE    = regexp.MustCompile(`(\w+)\s+(\w+)\s*=\s*(\w+)`)
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

	setTsTypeAndZeroValue()
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

func matchRequests(filename string, fileBytes []byte) {
	reMatched := RequestRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		// 没什么意义，因为匹配函数如果找到了匹配项，这里的长度就必然不小于3,但是因为后续要直接使用下标访问，所以还是判断一下，下同
		if len(reMatched[i]) < 3 {
			continue
		}

		requestName := string(reMatched[i][1])

		data.GeneratorIns.RequestAffiliation[filename] = append(data.GeneratorIns.RequestAffiliation[filename], requestName)
		data.GeneratorIns.Requests[requestName] = string(reMatched[i][2])
	}
}

func matchStructs(filename string, fileBytes []byte) {
	reMatched := StructRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 3 {
			continue
		}

		structureName := string(reMatched[i][1])

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], structureName)
		data.GeneratorIns.Structures[structureName] = &data.StructureItem{
			FromFile: filename,
			Typ:      &data.StructureType{IsStruct: true},
			Fields:   matchStructFields(reMatched[i][2]),
		}
	}
}

func matchStructFields(fieldBytes []byte) []*data.StructureField {
	fieldSlice := make([]*data.StructureField, 0)

	fieldSliceBytes := utils.BytesSplit(fieldBytes, '\n', ';')
	for _, v := range fieldSliceBytes {
		reMatched := StructFieldRE.FindSubmatch(v)
		if len(reMatched) < 3 {
			continue
		}

		typ := reMatched[1]
		isArray := bytes.HasPrefix(typ, []byte("[]"))
		typ = bytes.TrimPrefix(typ, []byte("[]"))
		typ = bytes.TrimPrefix(typ, []byte("*"))

		fieldIns := &data.StructureField{
			Name:       string(reMatched[2]),
			GoType:     string(typ),
			IsArray:    isArray,
			IsEmbedded: len(reMatched[2]) < 1,
		}

		fieldSlice = append(fieldSlice, fieldIns)
	}

	return fieldSlice
}

func matchEnums(filename string, fileBytes []byte) {
	reMatched := EnumRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 3 {
			continue
		}

		enumName := string(reMatched[i][1])

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], enumName)
		data.GeneratorIns.Structures[enumName] = &data.StructureItem{
			FromFile: filename,
			Typ:      &data.StructureType{IsEnum: true},
			Fields:   matchEnumUnits(enumName, reMatched[i][2]),
		}
	}
}

func matchEnumUnits(enumName string, fileBytes []byte) []*data.StructureField {
	enumUnitSlice := []*data.StructureField{{
		Name:        utils.EnumPlaceholder,
		GoType:      enumName,
		TSType:      enumName,
		TSZeroValue: "-1",
	}}

	reMatched := EnumUnitRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 4 {
			continue
		}

		enumName = string(reMatched[i][2])
		enumUnitIns := &data.StructureField{
			Name:        strings.TrimPrefix(string(reMatched[i][1]), enumName+"_"),
			GoType:      enumName,
			TSType:      enumName,
			TSZeroValue: string(reMatched[i][3]),
		}

		enumUnitSlice = append(enumUnitSlice, enumUnitIns)
	}

	return enumUnitSlice
}

// setTsTypeAndZeroValue set ts type and ts zero value for each field of go 'struct' typ
func setTsTypeAndZeroValue() {
	for _, structureIns := range data.GeneratorIns.Structures {
		if structureIns.Typ.IsEnum {
			continue // ignore enum types, enum has set ts type and ts zero value when matched
		}

		for i := 0; i < len(structureIns.Fields); i++ {
			fieldIns := structureIns.Fields[i]
			targetStructure, ok := data.GeneratorIns.Structures[fieldIns.GoType]

			if fieldIns.IsEmbedded {
				// 如果一个字段既不具名、也不带有json tag，则将其视作嵌入式结构体，
				// 移除该字段，将该字段所示嵌入式结构体的字段提高到与该字段同级
				structureIns.Fields = append(structureIns.Fields[:i], structureIns.Fields[i+1:]...) // rm embedded structure
				structureIns.Fields = append(structureIns.Fields, targetStructure.Fields...)        // add fields
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
				// 所以我们没有扫描到一个struct就设置一个struct的ts类型与零值，而是在全部扫描完成后统一设置
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
}
