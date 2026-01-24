// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.3

import { Pagination } from "./common.go"

export class CreateNoteReq {
    is_anonymous: boolean = false;
    title: string = "";
    content: string = "";
}

export class CreateNoteRes {}

export class Note {
    id: number = 0;
    created_at: number = 0;
    updated_at: number = 0;
    note_id: string = "";
    writer: string = ""; // writer nickname
    is_anonymous: boolean = false;
    title: string = "";
    content: string = "";
}

export class ListNoteReq {
    page: Pagination = new Pagination();
    list_my_flag: boolean = false;
}

export class ListNoteRes {
    amount: number = 0;
    notes: Array<Note> = new Array<Note>();
}

// ModifyNoteReq modify default is old value, only set fields not equal to old values
// can only modify myself note
export class ModifyNoteReq {
    id: number = 0;
    is_anonymous: boolean = false;
    title: string = "";
    content: string = "";
}

export class ModifyNoteRes {}

// DeleteNoteReq can only delete myself note
export class DeleteNoteReq {
    id: number = 0;
}

export class DeleteNoteRes {}
