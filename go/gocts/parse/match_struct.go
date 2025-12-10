package parse

import (
	"bytes"
	"regexp"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

var (
	StructRE      = regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	StructFieldRE = regexp.MustCompile(`\s*(?:\w+\s+)?([\[\]*\w]+)(?:\s+.*?json:["'](\w+)["'].*?)?`)
)

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
