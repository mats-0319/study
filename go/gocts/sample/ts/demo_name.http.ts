// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.1

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { ListUserRes, ListUserReq, UserIdentify, CreateUserRes } from "./demo_name.go"
import { Pagination } from "./common.go"

class DemoNameAxios {
    public listUser(operator: string, list_identify: UserIdentify, page: Pagination): Promise<AxiosResponse<ListUserRes>> {
        let req: ListUserReq = {
            operator: operator,
            list_identify: list_identify,
            page: page,
        }

        return axiosWrapper.post("/user/list", req)
    }

    public createUser(): Promise<AxiosResponse<CreateUserRes>> {
        return axiosWrapper.post("/user/create")
    }
}

export const demoNameAxios: DemoNameAxios = new DemoNameAxios()
