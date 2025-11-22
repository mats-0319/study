package generate_ts

import (
	"strconv"
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func GenerateConfig() {
	content := utils.Copyright
	content = append(content, formatConfigFile()...)

	utils.WriteFile(data.GeneratorIns.Config.TsDir+"config.ts", content)
}

func formatConfigFile() string {
	configStr := `
import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
{{ $indentation }}baseURL: "{{ $baseURL }}",
{{ $indentation }}timeout: {{ $timeout }},
});
`
	configStr = strings.ReplaceAll(configStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)
	configStr = strings.ReplaceAll(configStr, "{{ $baseURL }}", data.GeneratorIns.Config.BaseURL)
	configStr = strings.ReplaceAll(configStr, "{{ $timeout }}", strconv.Itoa(data.GeneratorIns.Config.Timeout))

	return configStr
}
