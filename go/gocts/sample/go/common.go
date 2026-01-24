package api

type Pagination struct {
	Size int `json:"size"`
	Num  int `json:"num"`
}

type Response struct {
	IsSuccess bool   `json:"is_success"`
	Err       string `json:"err"`
	Data      any    `json:"data"`
}
