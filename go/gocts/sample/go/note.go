package api

const URI_CreateNote = "/note/create"

type CreateNoteReq struct {
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type CreateNoteRes struct {
}

const URI_ListNote = "/note/list"

type Note struct {
	ID          uint   `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	NoteID      string `json:"note_id"`
	Writer      string `json:"writer"` // writer nickname
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type ListNoteReq struct {
	Page       Pagination `json:"page"`
	ListMyFlag bool       `json:"list_my_flag"`
}

type ListNoteRes struct {
	Amount int64   `json:"amount"`
	Notes  []*Note `json:"notes"`
}

const URI_ModifyNote = "/note/modify"

// ModifyNoteReq modify default is old value, only set fields not equal to old values
// can only modify myself note
type ModifyNoteReq struct {
	ID          uint   `json:"id"`
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type ModifyNoteRes struct {
}

const URI_DeleteNote = "/note/delete"

// DeleteNoteReq can only delete myself note
type DeleteNoteReq struct {
	ID uint `json:"id"`
}

type DeleteNoteRes struct {
}
