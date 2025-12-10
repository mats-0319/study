package gen

import (
	"fmt"
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func GenerateRequestFiles() {
	for filename := range data.GeneratorIns.RequestAffiliation {
		content := utils.Copyright
		content = append(content, formatRequestFile(filename)...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.RequestFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func formatRequestFile(filename string) string {
	requestFileStr := `
import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
{{ $importStructures }}

class {{ $filename_Big }}Axios {
{{ $requests }}
}

export const {{ $filename }}Axios: {{ $filename_Big }}Axios = new {{ $filename_Big }}Axios()
`

	externalStructures := make(map[string]*utils.Set)
	requestStr := formatHttpRequests(filename, externalStructures)

	importStructuresStr := formatStructuresImport(externalStructures)
	importStructuresStr = strings.TrimSpace(importStructuresStr)

	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $importStructures }}", importStructuresStr)
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $filename }}", filename)
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $filename_Big }}", utils.MustBig(filename))
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $requests }}", requestStr)
	requestFileStr = strings.ReplaceAll(requestFileStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return requestFileStr
}

// formatHttpRequests format http requests, also statistics external structures for declare their import statement
func formatHttpRequests(filename string, externalStructures map[string]*utils.Set) string {
	requestsStr := ""

	for _, requestName := range data.GeneratorIns.RequestAffiliation[filename] {
		structureItemIns := &data.StructureItem{}

		reqStructureName := requestName + data.GeneratorIns.Config.RequestStructureSuffix
		for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
			if structureName == reqStructureName { // find 'xxxReq' message
				structureItemIns, _ = data.GeneratorIns.Structures[structureName]
				break
			}
		}

		responseResName := requestName + data.GeneratorIns.Config.ResponseStructureSuffix
		externalStructures[filename] = externalStructures[filename].Add(responseResName)

		if len(structureItemIns.Fields) > 0 { // 'xxxReq' has fields
			externalStructures[filename] = externalStructures[filename].Add(reqStructureName)
		}

		requestsStr += formatOneHttpRequest(requestName, structureItemIns, externalStructures)
	}

	return strings.TrimSpace(requestsStr)
}

// formatOneHttpRequest
func formatOneHttpRequest(requestName string, structureItemIns *data.StructureItem, externalStructures map[string]*utils.Set) string {
	httpReqInvokeStr := `
{{ $indentation }}public {{ $requestName_Small }}({{ $inputParams }}): Promise<AxiosResponse<{{ $requestName }}{{ $responseStructureSuffix }}>> {{{ $requestParam }}
{{ $indentation }}{{ $indentation }}return axiosWrapper.post("{{ $requestURI }}"{{ $invokeParam }})
{{ $indentation }}}
`

	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestName_Small }}", utils.MustSmall(requestName))
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestName }}", requestName)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $responseStructureSuffix }}", data.GeneratorIns.Config.ResponseStructureSuffix)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestURI }}", data.GeneratorIns.Requests[requestName])

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
			if structureIns, ok := data.GeneratorIns.Structures[structureField.GoType]; ok {
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
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestMessageSuffix }}", data.GeneratorIns.Config.RequestStructureSuffix)

	return httpReqInputStr
}
