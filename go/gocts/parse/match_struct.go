package parse

import (
	"bytes"
	"regexp"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

var (
	StructRE      = regexp.MustCompile(`((?://.*?\n)+)?\s*type\s+(\w+)\s+struct\s*{([^}]*)}`)
	StructFieldRE = regexp.MustCompile(`\s*(?:\w+\s+)?([\[\]*\w]+)(?:\s+.*?json:["'](\w+)["'])?(?:.*?(//[^\n]+))?`)
)

func matchStructs(filename string, fileBytes []byte) {
	reMatched := StructRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		if len(reMatched[i]) < 4 {
			continue
		}

		structureName := string(reMatched[i][2])

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], structureName)
		data.GeneratorIns.Structures[structureName] = &data.StructureItem{
			Name:     structureName,
			FromFile: filename,
			Typ:      &data.StructureType{IsStruct: true},
			Comment:  string(reMatched[i][1]),
			Fields:   matchStructFields(reMatched[i][3]),
		}
	}
}

func matchStructFields(fieldBytes []byte) []*data.StructureField {
	fieldSlice := make([]*data.StructureField, 0)

	fieldSliceBytes := utils.BytesSplit(fieldBytes, '\n', ';')
	for _, v := range fieldSliceBytes {
		reMatched := StructFieldRE.FindSubmatch(v)
		if len(reMatched) < 4 {
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
			Comment:    string(reMatched[3]),
		}

		fieldSlice = append(fieldSlice, fieldIns)
	}

	return fieldSlice
}
