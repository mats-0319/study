package parse

import (
	"regexp"
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

var (
	EnumRE     = regexp.MustCompile(`type\s+(\w+)\s*=?\s*\w+\s+const\s*\(([^)]*)\)`)
	EnumUnitRE = regexp.MustCompile(`(\w+)\s+(\w+)\s*=\s*(\w+)`)
)

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
