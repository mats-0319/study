// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.1

import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
    baseURL: import.meta.env.Vite_axios_base_url,
    timeout: 3000,
});

// generate '.env.development' file in root path, 
// with content 'Vite_axios_base_url = "http://127.0.0.1:10319/api"'
//
// generate '.env.production' file in root path,
// with content 'Vite_axios_base_url = "https://xxx.xxx"'
//
// modify 'vite.config.ts', add config:
// envPrefix: "Vite_"
