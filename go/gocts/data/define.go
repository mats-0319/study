package data

type Generator struct {
	Config *GeneratorConfig

	RequestAffiliation   map[string][]string // filename - request name(s)
	StructureAffiliation map[string][]string // filename - structure name(s)

	Requests   map[string]string         // request name - request uri
	Structures map[string]*StructureItem // structure name - structure item

	TsBasicType map[string]string // go basic type - ts basic type
	TsZeroValue map[string]string // ts basic type - ts zero value

	IndentationStr string // indentation str
}

var GeneratorIns = &Generator{
	Config:               &GeneratorConfig{},
	RequestAffiliation:   make(map[string][]string),
	StructureAffiliation: make(map[string][]string),
	Requests:             make(map[string]string),
	Structures:           make(map[string]*StructureItem),
	TsBasicType:          make(map[string]string),
	TsZeroValue:          make(map[string]string),
}

type GeneratorConfig struct {
	// work path
	GoDir string `json:"go_dir"`
	TsDir string `json:"ts_dir"`
	// naming conventions 命名规范/约定
	RequestStructureSuffix  string `json:"request_structure_suffix"`
	ResponseStructureSuffix string `json:"response_structure_suffix"`
	RequestFileSuffix       string `json:"request_file_suffix"`
	StructureFileSuffix     string `json:"structure_file_suffix"`
	// params
	BasicGoType []struct {
		GoType      []string `json:"go_type"`
		TsType      string   `json:"ts_type"`
		TsZeroValue string   `json:"ts_zero_value"`
	} `json:"basic_go_type"`
	Indentation int `json:"indentation"`
}

type StructureItem struct {
	Name     string
	FromFile string
	Typ      *StructureType
	Comment  string
	Fields   []*StructureField
}

type StructureType struct {
	IsStruct bool
	IsEnum   bool
}

type StructureField struct {
	Name       string // field name, from json tag of go struct field
	GoType     string
	IsArray    bool
	IsEmbedded bool
	Comment    string

	TSType      string
	TSZeroValue string
}

var DefaultGeneratorConfig = &GeneratorConfig{
	GoDir:                   "./go/",
	TsDir:                   "./ts/",
	RequestStructureSuffix:  "Req",
	ResponseStructureSuffix: "Res",
	RequestFileSuffix:       ".http.ts",
	StructureFileSuffix:     ".go.ts",
	BasicGoType: []struct {
		GoType      []string `json:"go_type"`
		TsType      string   `json:"ts_type"`
		TsZeroValue string   `json:"ts_zero_value"`
	}{{
		GoType:      []string{"string", "byte", "rune"},
		TsType:      "string",
		TsZeroValue: `""`,
	}, {
		GoType: []string{"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64"},
		TsType:      "number",
		TsZeroValue: `0`,
	}, {
		GoType:      []string{"bool"},
		TsType:      "boolean",
		TsZeroValue: `false`,
	}, {
		GoType:      []string{"any", "interface{}"},
		TsType:      "Object",
		TsZeroValue: `{}`,
	}},
	Indentation: 4,
}
