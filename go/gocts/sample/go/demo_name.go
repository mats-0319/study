package api

type UserIdentify int8

const (
	UserIdentify_Second UserIdentify = 20
	UserIdentify_Value0 UserIdentify = 10
	UserIdentify_Value2 UserIdentify = 40
)

const URI_ListUser = "/user/list"

// ListUserReq comment
// some comment
type ListUserReq struct {
	//Operator     string       `json:"operator"` // this is a comment
	Operator     string       `json:"operator"` // this is a comment
	ListIdentify UserIdentify `json:"list_identify"`
	Pagination   `json:"page"`
}

type ListUserRes struct {
	*ResBase
	Summary int64 `json:"summary"`
	Users   any   `json:"users"`
}

const URI_CreateUser = "/user/create"

type CreateUserReq struct{}

type CreateUserRes struct {
	ResBase
}
