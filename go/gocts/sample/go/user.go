package api

const URI_Login = "/login"

type LoginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"` // hex(sha256('text'))
	TotpCode string `json:"totp_code"`
}

type LoginRes struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
}

const URI_Register = "/register"

type RegisterReq struct {
	UserName string `json:"user_name"` // nickname is same, user can modify later
	Password string `json:"password"`  // hex(sha256('text')), server generate 'salt' and save it
}

type RegisterRes struct {
}

const URI_ListUser = "/user/list"

type User struct {
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	UserName  string `json:"user_name"` // login name
	Nickname  string `json:"nickname"`  // display name
	TotpKey   string `json:"totp_key"`  // 允许为空，需要设置后启动
	IsAdmin   bool   `json:"is_admin"`
	LastLogin int64  `json:"last_login"` // timestamp, unit: milli
}

// ListUserReq 考虑添加按照字段查询、按照字段排序
type ListUserReq struct {
	Page Pagination `json:"page"`
}

type ListUserRes struct {
	Amount int64   `json:"amount"` // 符合查询条件的用户总数
	Users  []*User `json:"users"`
}

const URI_ModifyUser = "/user/modify"

// ModifyUserReq 属性字段为空，视为不修改对应字段，有专属的bool变量标识是否修改的字段不适用该默认规则
// 例如totp key字段，flag为true且key字段为空，视为禁用totp
type ModifyUserReq struct {
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`       // hex(sha256('text')), check: can't be same
	ModifyTkFlag bool   `json:"modify_tk_flag"` // if modify totp key
	TotpKey      string `json:"totp_key"`       // length: 16
}

type ModifyUserRes struct {
}

const URI_Authenticate = "/user/authenticate"

// AuthenticateReq 相当于使用token的登录
type AuthenticateReq struct {
}

type AuthenticateRes struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
}
