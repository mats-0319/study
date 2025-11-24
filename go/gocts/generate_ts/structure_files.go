package generate_ts

import (
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func GenerateStructureFiles() {
	for filename := range data.GeneratorIns.StructureAffiliation {
		content := utils.Copyright

		externalStructures := make(map[string][]string)
		structuresStr := ""
		for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
			structuresStr += formatStructure(structureName, externalStructures)
		}
		delete(externalStructures, filename) // not import from current file

		importStructuresStr := formatStructuresImport(externalStructures)

		content = append(content, importStructuresStr...)
		content = append(content, structuresStr...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.StructureFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func formatStructure(structureName string, externalStructures map[string][]string) string {
	structureItemIns, _ := data.GeneratorIns.Structures[structureName]

	structureStr := ""
	switch {
	case structureItemIns.Typ.IsStruct:
		structureStr = formatStruct(structureName, structureItemIns, externalStructures)
	case structureItemIns.Typ.IsEnum:
		structureStr = formatEnum(structureName, structureItemIns)
	}

	return structureStr
}

func formatStruct(structureName string, structureItemIns *data.StructureItem, externalStructures map[string][]string) string {
	fieldsStr := ""
	if len(structureItemIns.Fields) > 0 {
		fieldsStr = "\n"
	}
	for _, fieldIns := range structureItemIns.Fields {
		field := "{{ $indentation }}{{ $fieldName }}: {{ $fieldType_Ts }} = {{ $fieldZeroValue_Ts }};\n"
		field = strings.ReplaceAll(field, "{{ $fieldName }}", fieldIns.Name)
		field = strings.ReplaceAll(field, "{{ $fieldType_Ts }}", fieldIns.TSType)
		field = strings.ReplaceAll(field, "{{ $fieldZeroValue_Ts }}", fieldIns.TSZeroValue)

		fieldsStr += field

		if structureIns, ok := data.GeneratorIns.Structures[fieldIns.GoType]; ok {
			externalStructures[structureIns.FromFile] = append(externalStructures[structureIns.FromFile], fieldIns.TSType)
		}
	}

	structStr := "\nexport class {{ $structureName }} {{{ $structureFields }}}\n"
	structStr = strings.ReplaceAll(structStr, "{{ $structureName }}", structureName)
	structStr = strings.ReplaceAll(structStr, "{{ $structureFields }}", fieldsStr)
	structStr = strings.ReplaceAll(structStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return structStr
}

func formatEnum(enumName string, enumItemIns *data.StructureItem) string {
	enumUnitsStr := ""
	for _, enumUnitIns := range enumItemIns.Fields {
		unit := "{{ $indentation }}{{ $enumName }} = {{ $enumZeroValue_Ts }},\n"
		unit = strings.ReplaceAll(unit, "{{ $enumName }}", enumUnitIns.Name)
		unit = strings.ReplaceAll(unit, "{{ $enumZeroValue_Ts }}", enumUnitIns.TSZeroValue)

		enumUnitsStr += unit
	}

	enumStr := "\nexport enum {{ $enumName }} {\n{{ $enumUnits }}}\n"
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumName }}", enumName)
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumUnits }}", enumUnitsStr)
	enumStr = strings.ReplaceAll(enumStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return enumStr
}

// structures: from filename - structures' name
func formatStructuresImport(externalStructures map[string][]string) string {
	if len(externalStructures) < 1 {
		return ""
	}

	importStructuresStr := "\n"
	for fromFile, structureNames := range externalStructures {
		str := "import { {{ $structures }} } from \"./{{ $filename }}{{ $structureFileSuffix }}\"\n"
		str = strings.ReplaceAll(str, "{{ $structures }}", strings.Join(structureNames, ", "))
		str = strings.ReplaceAll(str, "{{ $filename }}", fromFile)

		importStructuresStr += str
	}
	importStructuresStr = strings.ReplaceAll(importStructuresStr, "{{ $structureFileSuffix }}", strings.TrimSuffix(data.GeneratorIns.Config.StructureFileSuffix, ".ts"))

	return importStructuresStr
}
