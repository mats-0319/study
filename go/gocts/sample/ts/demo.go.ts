// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.1

import { Pagination } from "./common.go"

// ListUserReq comment
// some comment
export class ListUserReq {
    operator: string = ""; // this is a comment
    list_identify: UserIdentify = UserIdentify.placeholder;
    page: Pagination = new Pagination();
}

export class ListUserRes {
    summary: number = 0;
    users: Object = {};
    is_success: boolean = false;
    err: string = "";
}

export class CreateUserReq {}

export class CreateUserRes {
    is_success: boolean = false;
    err: string = "";
}

export enum UserIdentify {
    placeholder = -1,
    Second = 20,
    Value0 = 10,
    Value2 = 40,
}
