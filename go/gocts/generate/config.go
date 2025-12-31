package gen

import (
	"strings"

	"github.com/mats9693/study/go/gocts/data"
	"github.com/mats9693/study/go/gocts/utils"
)

func GenerateConfigFile() {
	content := utils.Copyright
	content = append(content, formatConfigFile()...)

	utils.WriteFile(data.GeneratorIns.Config.TsDir+"config.ts", content)
}

func formatConfigFile() string {
	configStr := `
import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
	baseURL: getBaseUrl(),
	timeout: 3000
})

// getBaseUrl according to dev/product env & local ipv4 addr, set 'baseURL'
function getBaseUrl(): string {
	let url = import.meta.env.Vite_axios_base_url
	let localIP = window.location.hostname

	return import.meta.env.DEV ? url.replace("127.0.0.1", localIP) : url
}

// generate '.env.development' file in root path, 
// with content 'Vite_axios_base_url = "http://127.0.0.1:10319/api"'
//
// generate '.env.production' file in root path,
// with content 'Vite_axios_base_url = "https://xxx.xxx"'
//
// modify 'vite.config.ts', add config:
// envPrefix: "Vite_"
`
	configStr = strings.ReplaceAll(configStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return configStr
}
