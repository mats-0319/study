package token

import (
	"encoding/json"
	"log"
	"os"
)

type Generator struct {
	Config *GenConfig

	RequestFrom   map[string][]string // filename - request name(s)
	StructureFrom map[string][]string // filename - structure name(s)

	Requests   map[string]string         // request name - request uri
	Structures map[string]*StructureItem // structure name - structure item

	TsType map[string]*TsType // go basic type - ts type 'name' and 'zero value'

	IndentationStr string // indentation str
}

type StructureItem struct {
	Name     string
	FromFile string
	Typ      StructureType
	Comment  string
	Fields   []*StructureField
}

type StructureType int8

const (
	StructureType_Struct StructureType = 1
	StructureType_Enum   StructureType = 2
)

type StructureField struct {
	Name       string // field name, from json tag of go struct field
	GoType     string
	IsArray    bool
	IsEmbedded bool
	Comment    string

	// 这里的ts类型名称与零值，对应的是实际的字段类型。
	// 相比于ts基础类型，这里可能已经体现了数组
	TSType      string
	TSZeroValue string
}

type TsType struct {
	Name      string
	ZeroValue string
}

var GeneratorIns = &Generator{
	Config: &GenConfig{},

	RequestFrom:   make(map[string][]string),
	StructureFrom: make(map[string][]string),

	Requests:   make(map[string]string),
	Structures: make(map[string]*StructureItem),

	TsType: make(map[string]*TsType),
}

func (ins *Generator) Initialize(configFile string) {
	fileBytes, err := os.ReadFile(configFile)

	configIns := &GenConfig{}
	err2 := json.Unmarshal(fileBytes, configIns)

	if err != nil || err2 != nil {
		log.Println("> Gocts: Use default config, error: ", err, err2)
	} else {
		log.Println("> Gocts: Use config file: ", configFile)
	}

	configIns.mustValid()
	ins.Config = configIns

	// set type maps
	for _, typIns := range configIns.BasicGoType {
		tsBasicType := &TsType{
			Name:      typIns.TsType,
			ZeroValue: typIns.TsZeroValue,
		}

		for _, goTyp := range typIns.GoType {
			ins.TsType[goTyp] = tsBasicType
		}
	}

	// set indentation
	indentationBytes := make([]byte, 0, ins.Config.Indentation)
	for i := 0; i < ins.Config.Indentation; i++ {
		indentationBytes = append(indentationBytes, ' ')
	}
	ins.IndentationStr = string(indentationBytes)
}
