package models

import (
    "net/http"
    "github.com/enzdor/gomsg/sqlc"
)

type ErrorData struct {
	Status int
	Message string
}

var (
    NotFoundData = ErrorData{Status: http.StatusNotFound, Message: "Not found"}
    InternalServerErrorData = ErrorData{Status: http.StatusInternalServerError, Message: "Internal server error"}
)

type IndexData struct {
	Threads []sqlc.Thread
}

type BoardData struct {
	Threads []sqlc.Thread
	Name string
}

type ThreadData struct {
	Op sqlc.Thread
	Replies []sqlc.Reply
}

type PostData struct {
	Title string
	Comment string
	Board string
	Errors [2]FormError
}

type ReplyData struct {
	Comment string
	Thread_id int
	Error FormError
}

type KillData struct {
	Thread_id int
}
