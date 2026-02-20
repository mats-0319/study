// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.3

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { CreateNoteRes, CreateNoteReq, ListNoteRes, ListNoteReq, ModifyNoteRes, ModifyNoteReq, DeleteNoteRes, DeleteNoteReq } from "./note.go"
import { Pagination } from "./common.go"

class NoteAxios {
    public createNote(is_anonymous: boolean, title: string, content: string): Promise<AxiosResponse<CreateNoteRes>> {
        let req: CreateNoteReq = {
            is_anonymous: is_anonymous,
            title: title,
            content: content,
        }

        return axiosWrapper.post("/note/create", req)
    }

    public listNote(page: Pagination, list_my_flag: boolean): Promise<AxiosResponse<ListNoteRes>> {
        let req: ListNoteReq = {
            page: page,
            list_my_flag: list_my_flag,
        }

        return axiosWrapper.post("/note/list", req)
    }

    public modifyNote(id: number, is_anonymous: boolean, title: string, content: string): Promise<AxiosResponse<ModifyNoteRes>> {
        let req: ModifyNoteReq = {
            id: id,
            is_anonymous: is_anonymous,
            title: title,
            content: content,
        }

        return axiosWrapper.post("/note/modify", req)
    }

    public deleteNote(id: number): Promise<AxiosResponse<DeleteNoteRes>> {
        let req: DeleteNoteReq = {
            id: id,
        }

        return axiosWrapper.post("/note/delete", req)
    }
}

export const noteAxios: NoteAxios = new NoteAxios()
