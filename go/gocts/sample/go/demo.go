package api

type UserIdentify int8

const (
	UserIdentify_Second UserIdentify = 20
	UserIdentify_Value0 UserIdentify = 10
	UserIdentify_Value2 UserIdentify = 40
)

const URI_ListUser = "/user/list"

type ListUserReq struct {
	Operator     string       `json:"operator"`
	ListIdentify UserIdentify `json:"list_identify"`
	Page         Pagination   `json:"page"`
}

type ListUserRes struct {
	ResBase `json:"res"`
	Summary int64    `json:"summary"`
	Users   []string `json:"users"`
}

const URI_CreateUser = "/user/create"

type CreateUserReq struct{}

type CreateUserRes struct {
	ResBase
}
