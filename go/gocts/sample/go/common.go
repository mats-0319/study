package api

type ResBase struct {
	IsSuccess bool   `json:"is_success"`
	Err       string `json:"err"`
}

type Pagination struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type ResType interface {
	string | ListUserRes | CreateUserRes
}
