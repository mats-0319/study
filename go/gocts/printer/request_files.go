package gen

import (
	"fmt"
	"strings"

	"github.com/mats0319/study/go/gocts/token"
	"github.com/mats0319/study/go/gocts/utils"
)

func GenerateRequestFiles() {
	for fileName := range token.GeneratorIns.RequestFrom {
		content := utils.Copyright
		content = append(content, formatRequestFile(fileName)...)

		absolutePath := token.GeneratorIns.Config.TsDir + fileName + token.GeneratorIns.Config.RequestFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func formatRequestFile(fileName string) string {
	requestFileStr := `
import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
{{ $importStructures }}

class {{ $filename_Big }}Axios {
{{ $requests }}
}

export const {{ $filename_Small }}Axios: {{ $filename_Big }}Axios = new {{ $filename_Big }}Axios()
`

	externalStructures := make(map[string]*utils.Set)
	requestStr := formatHttpRequests(fileName, externalStructures)

	importStructuresStr := formatStructuresImport(externalStructures)
	importStructuresStr = strings.TrimSpace(importStructuresStr)

	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $importStructures }}", importStructuresStr)
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $filename_Small }}", utils.MustSmall(fileName))
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $filename_Big }}", utils.MustBig(fileName))
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $requests }}", requestStr)
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $indentation }}", token.GeneratorIns.IndentationStr)

	return requestFileStr
}

// formatHttpRequests format http requests, also statistics external structures for declare their import statement
func formatHttpRequests(fileName string, externalStructures map[string]*utils.Set) string {
	var requestsStr strings.Builder

	for _, requestName := range token.GeneratorIns.RequestFrom[fileName] {
		structureItemIns := &token.StructureItem{}

		reqStructureName := requestName + token.GeneratorIns.Config.RequestStructureSuffix
		for _, structureName := range token.GeneratorIns.StructureFrom[fileName] {
			if structureName == reqStructureName { // find 'xxxReq' message
				structureItemIns, _ = token.GeneratorIns.Structures[structureName]
				break
			}
		}

		responseResName := requestName + token.GeneratorIns.Config.ResponseStructureSuffix
		externalStructures[fileName] = externalStructures[fileName].Add(responseResName)

		if len(structureItemIns.Fields) > 0 { // 'xxxReq' has fields
			externalStructures[fileName] = externalStructures[fileName].Add(reqStructureName)
		}

		requestsStr.WriteString(formatOneHttpRequest(requestName, structureItemIns, externalStructures))
	}

	return strings.TrimSpace(requestsStr.String())
}

// formatOneHttpRequest
func formatOneHttpRequest(requestName string, structureItemIns *token.StructureItem, externalStructures map[string]*utils.Set) string {
	httpReqInvokeStr := `
{{ $indentation }}public {{ $requestName_Small }}({{ $inputParams }}): Promise<AxiosResponse<{{ $requestName }}{{ $responseStructureSuffix }}>> {{{ $requestParam }}
{{ $indentation }}{{ $indentation }}return axiosWrapper.post("{{ $requestURI }}"{{ $invokeParam }})
{{ $indentation }}}
`

	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestName_Small }}", utils.MustSmall(requestName))
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestName }}", requestName)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $responseStructureSuffix }}", token.GeneratorIns.Config.ResponseStructureSuffix)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestURI }}", token.GeneratorIns.Requests[requestName])

	var (
		inputParamStr   string
		requestParamStr string
		invokeParamStr  string
	)
	if len(structureItemIns.Fields) > 0 { // 'xxxReq' has fields
		inputParamsSlice := make([]string, 0, len(structureItemIns.Fields))
		requestParamsSlice := make([]string, 0, len(structureItemIns.Fields))
		for _, structureField := range structureItemIns.Fields {
			inputParamsSlice = append(inputParamsSlice, fmt.Sprintf("%s: %s", structureField.Name, structureField.TSType))
			requestParamsSlice = append(requestParamsSlice, fmt.Sprintf("%s: %s,", structureField.Name, structureField.Name))

			// record if structure field's type need import from another file
			if structureIns, ok := token.GeneratorIns.Structures[structureField.GoType]; ok {
				externalStructures[structureIns.FromFile] = externalStructures[structureIns.FromFile].Add(structureField.TSType)
			}
		}

		inputParamStr = strings.Join(inputParamsSlice, ", ")
		requestParamStr = formatHttpRequestParam(requestName, requestParamsSlice)
		invokeParamStr = ", req"
	}

	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $inputParams }}", inputParamStr)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestParam }}", requestParamStr)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $invokeParam }}", invokeParamStr)

	return httpReqInvokeStr
}

func formatHttpRequestParam(requestName string, requestInputs []string) string {
	fieldStr := ""
	for i := range requestInputs {
		field := "{{ $indentation }}{{ $indentation }}{{ $indentation }}{{ $requestInput }}\n"
		field = strings.ReplaceAll(field, "{{ $requestInput }}", requestInputs[i])

		fieldStr += field
	}
	fieldStr = strings.TrimSuffix(fieldStr, "\n")

	httpReqInputStr := `
{{ $indentation }}{{ $indentation }}let req: {{ $requestName }}{{ $requestMessageSuffix }} = {
{{ $requestFields }}
{{ $indentation }}{{ $indentation }}}
`
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestFields }}", fieldStr)
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestName }}", requestName)
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestMessageSuffix }}", token.GeneratorIns.Config.RequestStructureSuffix)

	return httpReqInputStr
}
