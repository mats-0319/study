// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.3

import { Pagination } from "./common.go"

export class LoginReq {
    user_name: string = "";
    password: string = ""; // hex(sha256('text'))
    totp_code: string = "";
}

export class LoginRes {
    user_id: number = 0;
    user_name: string = "";
    nickname: string = "";
    is_admin: boolean = false;
}

export class RegisterReq {
    user_name: string = ""; // nickname is same, user can modify later
    password: string = ""; // hex(sha256('text')), server generate 'salt' and save it
}

export class RegisterRes {}

export class User {
    id: number = 0;
    created_at: number = 0;
    updated_at: number = 0;
    user_name: string = ""; // login name
    nickname: string = ""; // display name
    totp_key: string = ""; // 允许为空，需要设置后启动
    is_admin: boolean = false;
    last_login: number = 0; // timestamp, unit: milli
}

// ListUserReq 考虑添加按照字段查询、按照字段排序
export class ListUserReq {
    page: Pagination = new Pagination();
}

export class ListUserRes {
    amount: number = 0; // 符合查询条件的用户总数
    users: Array<User> = new Array<User>();
}

// ModifyUserReq 属性字段为空，视为不修改对应字段，有专属的bool变量标识是否修改的字段不适用该默认规则
// 例如totp key字段，flag为true且key字段为空，视为禁用totp
export class ModifyUserReq {
    nickname: string = "";
    password: string = ""; // hex(sha256('text')), check: can't be same
    modify_tk_flag: boolean = false; // if modify totp key
    totp_key: string = ""; // length: 16
}

export class ModifyUserRes {}

// AuthenticateReq 相当于使用token的登录
export class AuthenticateReq {}

export class AuthenticateRes {
    user_id: number = 0;
    user_name: string = "";
    nickname: string = "";
    is_admin: boolean = false;
}
