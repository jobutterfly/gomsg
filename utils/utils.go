package utils

import (
    "strconv"
    "net/http"
    "context"
    "strings"
    "html/template"
    "path/filepath"
    "log"

    "github.com/enzdor/gomsg/models"
    "github.com/enzdor/gomsg/sqlc"
)

func Serve(page string) *template.Template {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", page + ".html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
	    log.Fatal(err)
	}

	return tmpl
}

func ValidatePost(title string, comment string) ([2]models.FormError, error) {
    errors := [2]models.FormError{
	{
	    Bool: false,
	    Message: "",
	    Field: "title",
	},
	{
	    Bool: false,
	    Message: "",
	    Field: "comment",
	},
    }

    if strings.TrimSpace(title) == "" {
	errors[0] = models.FormError{
	    Bool: true,
	    Message: "This field is required",
	    Field: "title",
	}

    }

    if strings.TrimSpace(comment) == "" {
	errors[1] = models.FormError{
	    Bool: true,
	    Message: "This field is required",
	    Field: "comment",
	}

    }

    if errors[0].Bool || errors[1].Bool {
	err := &models.ValidateError{Message: "One of the fields has not passed the required validation rules."}
	return errors, err
    }

    return errors, nil
}

func ValidateReply(title string) (models.FormError, error) {
    error := models.FormError{
	Bool: false,
	Message: "",
	Field: "comment",
    }

    if strings.TrimSpace(title) == "" {
	error = models.FormError{
	    Bool: true,
	    Message: "This field is required",
	    Field: "comment",
	}

    }

    if error.Bool {
	err := &models.ValidateError{Message: "The field has not passed the required validation rules."}
	return error, err
    }

    return error, nil
}


func GetBoardData(queries *sqlc.Queries, id int32, name string) (models.BoardData, error){
    threads, err := queries.GetBoardThreads(context.Background(), id)
    if err != nil {
	data := models.BoardData{
	    Threads: []sqlc.Thread{},
	    Name: name,
	}
	return data, err
    }

    data := models.BoardData{
	Threads: threads,
	Name: name,
    }

    return data, nil
}

func GetThreadData(queries *sqlc.Queries, id int32) (models.ThreadData, error){
    errdata := models.ThreadData{
	Op: sqlc.Thread{},
	Replies: []sqlc.Reply{},
    }

    thread, err := queries.GetThread(context.Background(), id)
    if err != nil {
	return errdata, err
    }

    replies, err := queries.GetThreadReplies(context.Background(), id)
    if err != nil {
	return errdata, err
    }

    data := models.ThreadData{
	Op: thread,
	Replies: replies,
    }

    return data, nil
}

func CreateErrorData(status int) models.ErrorData{
    switch status {
    case http.StatusNotFound:
	return models.ErrorData{
	    Status: status,
	    Message: "Not found",
	}
    default:
	return models.ErrorData{
	    Status: http.StatusInternalServerError,
	    Message: "Internal server error",
	}
    }
}

type PathInfo struct {
    BoardName string
    ThreadID int
    Status int
}

type PathWant struct {
    BoardName bool
    ThreadID bool
    Status bool

}

func GetPathValues(ps []string, pw PathWant) (PathInfo, error){
    r := PathInfo{
	BoardName: "",
	ThreadID: 0,
	Status: 0,
    }

    if len(ps) > 3 {
	if ps[3] != "" {
	    err := &models.PathError{Message: "Not found"}
	    return r, err
	}
    }

    if pw.BoardName == true {
	name := ps[2]
	if name != "sports" && name != "random" && name != "tech" {
	    err := &models.PathError{Message: "Not found"}
	    return r, err
	}
	r.BoardName = name
	return r, nil
    } else if pw.ThreadID == true {
	thread_id, err := strconv.Atoi(ps[2])
	if err != nil {
	    return r, err
	}
	r.ThreadID = thread_id
	return r, nil
    } else if pw.Status == true {
	status, err := strconv.Atoi(ps[2])
	if err != nil {
	    return r, err
	}
	r.Status = status
	return r, nil
    }

    err := &models.PathError{Message: "Not found"}
    return r, err
}

func GetBoardID(name string) int32{
    var id int32

    switch name {
    case "sports":
	id = 1;
    case "random":
	id = 2;
    case "tech":
	id = 3;
    }

    return id
}













