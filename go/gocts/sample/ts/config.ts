// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v1.0.0

import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
    baseURL: "http://127.0.0.1:10319/api",
    timeout: 3000,
});
