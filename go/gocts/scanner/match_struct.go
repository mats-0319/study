package scanner

import (
	"bytes"
	"regexp"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

var (
	structRE      = regexp.MustCompile(`((?://.*?\n)+)?\s*type\s+(\w+)\s+struct\s*{([^}]*)}`)
	structFieldRE = regexp.MustCompile(`\s*(?:\w+\s+)?([\[\]*\w]+)(?:\s+.*?json:["'](\w+)["'])?(?:.*?(//[^\n]+))?`)
)

func matchStructs(fileName string, fileBytes []byte) {
	reMatched := structRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 4 {
			continue
		}

		structureName := string(reMatched[i][2])

		token.GeneratorIns.StructureFrom[fileName] = append(token.GeneratorIns.StructureFrom[fileName], structureName)
		token.GeneratorIns.Structures[structureName] = &token.StructureItem{
			Name:     structureName,
			FromFile: fileName,
			Typ:      token.StructureType_Struct,
			Comment:  string(reMatched[i][1]),
			Fields:   matchStructFields(reMatched[i][3]),
		}
	}
}

func matchStructFields(allFields []byte) []*token.StructureField {
	fieldSlice := make([]*token.StructureField, 0)

	fieldSliceSplit := utils.BytesSplit(allFields, '\n', ';')
	for _, fieldBytes := range fieldSliceSplit {
		reMatched := structFieldRE.FindSubmatch(fieldBytes)
		if len(reMatched) < 4 {
			continue
		}

		typ := reMatched[1]
		isArray := bytes.HasPrefix(typ, []byte("[]"))
		typ = bytes.TrimPrefix(typ, []byte("[]"))
		typ = bytes.TrimPrefix(typ, []byte("*"))

		fieldIns := &token.StructureField{
			Name:       string(reMatched[2]),
			GoType:     string(typ),
			IsArray:    isArray,
			IsEmbedded: len(reMatched[2]) < 1,
			Comment:    string(reMatched[3]),
		}

		fieldSlice = append(fieldSlice, fieldIns)
	}

	return fieldSlice
}
