package gen

import (
	"strings"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

func GenerateStructureFiles() {
	for fileName := range token.GeneratorIns.StructureFrom {
		content := utils.Copyright

		externalStructures := make(map[string]*utils.Set)
		var structuresStr strings.Builder
		for _, structureName := range token.GeneratorIns.StructureFrom[fileName] {
			structuresStr.WriteString(formatStructure(structureName, externalStructures))
		}
		delete(externalStructures, fileName) // not import from current file

		importStructuresStr := formatStructuresImport(externalStructures)

		content = append(content, importStructuresStr...)
		content = append(content, structuresStr.String()...)

		absolutePath := token.GeneratorIns.Config.TsDir + fileName + token.GeneratorIns.Config.StructureFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func formatStructure(structureName string, externalStructures map[string]*utils.Set) string {
	structureItemIns, _ := token.GeneratorIns.Structures[structureName]

	structureStr := ""
	switch {
	case structureItemIns.Typ == token.StructureType_Struct:
		structureStr = formatStruct(structureName, structureItemIns, externalStructures)
	case structureItemIns.Typ == token.StructureType_Enum:
		structureStr = formatEnum(structureName, structureItemIns)
	}

	return structureStr
}

func formatStruct(structureName string, structureItemIns *token.StructureItem, externalStructures map[string]*utils.Set) string {
	fieldsStr := ""
	if len(structureItemIns.Fields) > 0 {
		fieldsStr = "\n"
	}
	for _, fieldIns := range structureItemIns.Fields {
		field := "{{ $indentation }}{{ $fieldName }}: {{ $fieldType_Ts }} = {{ $fieldZeroValue_Ts }}{{ $fieldComment }}\n"
		field = strings.ReplaceAll(field, "{{ $fieldName }}", fieldIns.Name)
		field = strings.ReplaceAll(field, "{{ $fieldType_Ts }}", fieldIns.TSType)
		field = strings.ReplaceAll(field, "{{ $fieldZeroValue_Ts }}", fieldIns.TSZeroValue)
		comment := ""
		if len(fieldIns.Comment) > 0 {
			comment = " " + fieldIns.Comment
		}
		field = strings.ReplaceAll(field, "{{ $fieldComment }}", comment)

		fieldsStr += field

		if structureIns, ok := token.GeneratorIns.Structures[fieldIns.GoType]; ok {
			externalStructures[structureIns.FromFile] = externalStructures[structureIns.FromFile].Add(fieldIns.TSType)
		}
	}

	structStr := "\n{{ $structComment }}export class {{ $structName }} {{{ $structFields }}}\n"
	structStr = strings.ReplaceAll(structStr, "{{ $structComment }}", structureItemIns.Comment)
	structStr = strings.ReplaceAll(structStr, "{{ $structName }}", structureName)
	structStr = strings.ReplaceAll(structStr, "{{ $structFields }}", fieldsStr)
	structStr = strings.ReplaceAll(structStr, "{{ $indentation }}", token.GeneratorIns.IndentationStr)

	return structStr
}

func formatEnum(enumName string, enumItemIns *token.StructureItem) string {
	var enumUnitsStr strings.Builder
	for _, enumUnitIns := range enumItemIns.Fields {
		unit := "{{ $indentation }}{{ $enumName }} = {{ $enumZeroValue_Ts }},\n"
		unit = strings.ReplaceAll(unit, "{{ $enumName }}", enumUnitIns.Name)
		unit = strings.ReplaceAll(unit, "{{ $enumZeroValue_Ts }}", enumUnitIns.TSZeroValue)

		enumUnitsStr.WriteString(unit)
	}

	enumStr := "\nexport enum {{ $enumName }} {\n{{ $enumUnits }}}\n"
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumName }}", enumName)
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumUnits }}", enumUnitsStr.String())
	enumStr = strings.ReplaceAll(enumStr, "{{ $indentation }}", token.GeneratorIns.IndentationStr)

	return enumStr
}

// structures: from filename - structures' name
func formatStructuresImport(externalStructures map[string]*utils.Set) string {
	if len(externalStructures) < 1 {
		return ""
	}

	importStructuresStr := "\n"
	for fromFile, structureNames := range externalStructures {
		str := "import { {{ $structures }} } from \"./{{ $filename }}{{ $structureFileSuffix }}\"\n"
		str = strings.ReplaceAll(str, "{{ $structures }}", strings.Join(structureNames.Data, ", "))
		str = strings.ReplaceAll(str, "{{ $filename }}", fromFile)

		importStructuresStr += str
	}
	importStructuresStr = strings.ReplaceAll(importStructuresStr, "{{ $structureFileSuffix }}", strings.TrimSuffix(token.GeneratorIns.Config.StructureFileSuffix, ".ts"))

	return importStructuresStr
}
