package token

import "github.com/mats0319/study/go/gocts/utils"

type GenConfig struct {
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

var DefaultGeneratorConfig = &GenConfig{
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

// make sure all configs are valid, use default config cover empty ones
func (c *GenConfig) mustValid() {
	if len(c.GoDir) < 1 {
		c.GoDir = DefaultGeneratorConfig.GoDir
	} else {
		c.GoDir = utils.MustSuffix(c.GoDir, "/")
	}
	if len(c.TsDir) < 1 {
		c.TsDir = DefaultGeneratorConfig.TsDir
	} else {
		c.TsDir = utils.MustSuffix(c.TsDir, "/")
	}
	utils.MustExistDir(c.GoDir)
	utils.MustExistDir(c.GoDir + utils.GoBackupFolderName + "/")
	utils.EmptyDir(c.TsDir)

	if len(c.RequestStructureSuffix) < 1 {
		c.RequestStructureSuffix = DefaultGeneratorConfig.RequestStructureSuffix
	}
	if len(c.ResponseStructureSuffix) < 1 {
		c.ResponseStructureSuffix = DefaultGeneratorConfig.ResponseStructureSuffix
	}
	if len(c.RequestFileSuffix) < 1 {
		c.RequestFileSuffix = DefaultGeneratorConfig.RequestFileSuffix
	} else {
		c.RequestFileSuffix = utils.MustSuffix(c.RequestFileSuffix, ".ts")
	}
	if len(c.StructureFileSuffix) < 1 {
		c.StructureFileSuffix = DefaultGeneratorConfig.StructureFileSuffix
	} else {
		c.StructureFileSuffix = utils.MustSuffix(c.StructureFileSuffix, ".ts")
	}

	if len(c.BasicGoType) < 1 {
		c.BasicGoType = DefaultGeneratorConfig.BasicGoType
	}
	if c.Indentation < 1 {
		c.Indentation = DefaultGeneratorConfig.Indentation
	}
}
