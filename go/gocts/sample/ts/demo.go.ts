// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v1.0.0

import { Pagination, ResBase } from "./common.go"

export class ListUserReq {
    operator: string = "";
    list_identify: UserIdentify = UserIdentify.placeholder;
    page: Pagination = new Pagination();
}

export class ListUserRes {
    res: ResBase = new ResBase();
    summary: number = 0;
    users: Array<string> = new Array<string>();
}

export class CreateUserReq {}

export class CreateUserRes {}

export enum UserIdentify {
    placeholder = -1,
    Second = 20,
    Value0 = 10,
    Value2 = 40,
}
