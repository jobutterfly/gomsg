package controllers

import (
	"time"
	"strings"
	"context"
	"strconv"
	"net/http"
	"log"

	"github.com/enzdor/gomsg/models"
	"github.com/enzdor/gomsg/utils"
	"github.com/enzdor/gomsg/sqlc"
)

func (h *Handler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := utils.Serve("index")
	ps := strings.Split(r.URL.Path, "/")

	if len(ps) > 1 {
	    if ps[1] != "" {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
		return
	    }
	}

	threads, err := h.q.GetThreads(context.Background(), 3)
	if err != nil {
	    log.Fatal(err)
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
	    return
	}

	data := models.IndexData{
	    Threads: threads,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}


func (h *Handler) ServeBoard(w http.ResponseWriter, r *http.Request) {
	tmpl := utils.Serve("board")
	
	vs, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), utils.PathWant{BoardName: true, ThreadID: false, Status: false})
	if err != nil {
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
	    return
	}
	name := vs.BoardName
	id := utils.GetBoardID(name)

	data, err := utils.GetBoardData(h.q, id, name)
	if err != nil {
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
	    return
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}


func (h *Handler) ServeThread(w http.ResponseWriter, r *http.Request) {
	tmpl := utils.Serve("thread")

	vs, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), utils.PathWant{BoardName: false, ThreadID: true, Status: false})
	if err != nil {
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
	    return
	}
	id := vs.ThreadID

	data, err := utils.GetThreadData(h.q, int32(id))
	if err != nil {
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
	    return
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) ServePost(w http.ResponseWriter, r *http.Request) {
	tmpl := utils.Serve("post")
	method := r.Method

	vs, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), utils.PathWant{BoardName: true, ThreadID: false, Status: false})
	if err != nil {
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
	    return
	}
	name := vs.BoardName

	switch method {
	case "GET":
	    data := models.PostData{
		Title: "",
		Comment: "",
		Board: name,
		Errors: [2]models.FormError{
		    {Bool: false, Message: "", Field: "title"},
		    {Bool: false, Message: "", Field: "comment"},
		},
	    }
	    tmpl.ExecuteTemplate(w, "layout", data)
	    return

	case "POST":
	    id := utils.GetBoardID(name)

	    if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
		return
	    }

	    errors, err := utils.ValidatePost(r.FormValue("title"), r.FormValue("comment"))
	    if err != nil {
		data := models.PostData{
		    Title: r.FormValue("title"),
		    Comment: r.FormValue("comment"),
		    Board: name,
		    Errors: errors,
		}
		data.Errors = errors
		tmpl.ExecuteTemplate(w, "layout", data)
		return
	    }

	    nr, err := h.q.CountThreads(context.Background()); if err != nil {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
		return
	    }

	    if nr >= 20 {
		oldestThread, err := h.q.GetOldestThread(context.Background(), id)
		if err != nil {
		    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
		    return
		}
		_, nerr := h.q.DeleteThread(context.Background(), int32(oldestThread.ThreadID))
		if nerr != nil{
		    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
		    return
		}
	    }

	    if _, err := h.q.CreateThread(context.Background(), sqlc.CreateThreadParams{
		Title: r.FormValue("title"),
		Comment: r.FormValue("comment"),
		Date: strconv.Itoa(int(time.Now().Unix())),
		BoardID: id,
	    }); err != nil {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
		return
	    }

	    http.Redirect(w, r, "/board/" + name, http.StatusSeeOther)
	    return
	}

}

func (h *Handler) ServeReply(w http.ResponseWriter, r *http.Request) {
	tmpl := utils.Serve("reply")
	method := r.Method

	vs, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), utils.PathWant{BoardName: false, ThreadID: true, Status: false})
	if err != nil {
	    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
	    return
	}
	id := vs.ThreadID

	switch method {
	case "GET":
	    data := models.ReplyData{
		Comment: "",
		Thread_id: id,
		Error: models.FormError{
		    Bool: false, 
		    Message: "", 
		    Field: "",
		},
	    }
	    tmpl.ExecuteTemplate(w, "layout", data)
	    return
	case "POST":
	    if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
		return
	    }

	    error, err := utils.ValidateReply(r.FormValue("comment"))
	    if err != nil {
		data := models.ReplyData{
		    Comment: r.FormValue("comment"),
		    Thread_id: id,
		    Error: error,
		}
		tmpl.ExecuteTemplate(w, "layout", data)
		return
	    }

	    nr, err := h.q.CountReplies(context.Background(), int32(id)); if err != nil {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
		return
	    }

	    if nr >= 20 {
		_, err := h.q.DeleteThread(context.Background(), int32(id)); if err != nil{
		    http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
		    return
		}
		http.Redirect(w, r, "/kill/" + strconv.Itoa(id), http.StatusSeeOther)
		return
	    }

	    if _, err := h.q.CreateReply(context.Background(), sqlc.CreateReplyParams{
		Comment: r.FormValue("comment"),
		Date: strconv.Itoa(int(time.Now().Unix())),
		ThreadID: int32(id),
	    }); err != nil {
		http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusInternalServerError), http.StatusSeeOther)
		return
	    }

	    http.Redirect(w, r, "/thread/" + strconv.Itoa(id), http.StatusSeeOther)
	    return
	}
}

func (h *Handler) ServeKill(w http.ResponseWriter, r *http.Request) {
    tmpl := utils.Serve("kill")

    vs, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), utils.PathWant{BoardName: false, ThreadID: true, Status: false})
    if err != nil {
	http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
	return
    }
    id := vs.ThreadID

    data := models.KillData{
	Thread_id: id,
    }

    tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) ServeError(w http.ResponseWriter, r *http.Request) {
    tmpl := utils.Serve("error")

    vs, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), utils.PathWant{BoardName: false, ThreadID: false, Status: true})
    if err != nil {
	http.Redirect(w, r, "/error/" + strconv.Itoa(http.StatusNotFound), http.StatusSeeOther)
	return
    }
    status := vs.Status

    data := utils.CreateErrorData(status)

    tmpl.ExecuteTemplate(w, "layout", data)
}











