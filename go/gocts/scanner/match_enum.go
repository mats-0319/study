package scanner

import (
	"regexp"
	"strings"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

var (
	enumRE     = regexp.MustCompile(`type\s+(\w+)\s*=?\s*\w+\s+const\s*\(([^)]*)\)`)
	enumUnitRE = regexp.MustCompile(`(\w+)\s+(?:\w+)?\s*=\s*(-?\d+)`)
)

func matchEnums(fileName string, fileBytes []byte) {
	reMatched := enumRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 3 {
			continue
		}

		enumName := string(reMatched[i][1])

		token.GeneratorIns.StructureFrom[fileName] = append(token.GeneratorIns.StructureFrom[fileName], enumName)
		token.GeneratorIns.Structures[enumName] = &token.StructureItem{
			FromFile: fileName,
			Typ:      token.StructureType_Enum,
			Fields:   matchEnumUnits(enumName, reMatched[i][2]),
		}
	}
}

func matchEnumUnits(enumName string, fileBytes []byte) []*token.StructureField {
	enumUnitSlice := []*token.StructureField{{
		Name:        utils.EnumPlaceholder,
		GoType:      enumName,
		TSType:      enumName,
		TSZeroValue: "-1",
	}}

	reMatched := enumUnitRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 3 {
			continue
		}

		enumUnitIns := &token.StructureField{
			Name:        strings.TrimPrefix(string(reMatched[i][1]), enumName+"_"),
			GoType:      enumName,
			TSType:      enumName,
			TSZeroValue: string(reMatched[i][2]),
		}

		enumUnitSlice = append(enumUnitSlice, enumUnitIns)
	}

	return enumUnitSlice
}
