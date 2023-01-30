package controllers

import (
	"testing"
	"bytes"
	"strings"
	"io"
	"io/ioutil"
	"os"
	"context"
	"strconv"
	"time"
	"html/template"
	"net/http"
	"net/http/httptest"

	_ "github.com/go-sql-driver/mysql"
	"github.com/enzdor/gomsg/models"
	"github.com/enzdor/gomsg/utils"
	"github.com/enzdor/gomsg/sqlc"
	"github.com/joho/godotenv"
)

var Th *Handler

func stringTemplate(tmpl *template.Template, data any) (string, error){
    var buff bytes.Buffer 
    if err := tmpl.ExecuteTemplate(&buff, "layout", data); err != nil {
	return "", err
    }

    return buff.String(), nil
}

func start() error{
    if err := godotenv.Load("../.env"); err != nil {
	return err
    }
    user := os.Getenv("DBUSER")
    pass := os.Getenv("DBPASS")
    name := os.Getenv("TESTDBNAME")

    db := NewDB(user, pass, name)

    if _, err := db.Query("DELETE FROM replies; "); err != nil {
	return err
    }
    if _, err := db.Query("DELETE FROM threads; "); err != nil {
	return err
    }

    Th = NewHandler(db)

    return nil
}

func testRedirect(path string, h func(w http.ResponseWriter, r *http.Request)) error{
    redirectTestCase := struct {
	req 	*http.Request
	w	*httptest.ResponseRecorder
    } {
	req: httptest.NewRequest(http.MethodGet, path, nil),
	w: httptest.NewRecorder(), 
    }

    Th.ServeIndex(redirectTestCase.w, redirectTestCase.req)
    res := redirectTestCase.w.Result()
    defer res.Body.Close()

    url, err := res.Location()
    if err != nil {
	return err
    }

    if url.Path != "/error/404" {
	return &models.PathError{Message: "expected path to be /error/404 but got" + url.Path}
    }

    return nil
}

func cleanString (s string) string{
    s = strings.ReplaceAll(s, "\n", "")
    s = strings.ReplaceAll(s, "\t", "")
    s = strings.ReplaceAll(s, " ", "")

    return s
}

func TestServeIndex(t *testing.T){
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    createThreads := []sqlc.CreateThreadParams{
	{
	    Title: "This is the first title",
	    Comment: "This is the first comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 1,
	},
	{
	    Title: "This is the second title",
	    Comment: "This is the second comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 1,
	},
	{
	    Title: "This is the third title",
	    Comment: "This is the third comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 3,
	},
    }

    // populating db with threads that are going to be queried

    for _, tt := range createThreads {
	_, err := Th.q.CreateThread(context.Background(), tt)
	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
    }

    threads, err := Th.q.GetThreads(context.Background(), 3)
    if err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    testCases := []struct {
	name	string
	req 	*http.Request
	w	*httptest.ResponseRecorder
	te	*template.Template
	te_data	any
    }{
	{
	    name: "index",
	    req: httptest.NewRequest(http.MethodGet, "/", nil),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("index"),
	    te_data: models.IndexData{
		Threads: threads,
	    },
	},
    } 

    for _, tc := range testCases {
	t.Run(tc.name, func(t *testing.T){
	    ts , err := stringTemplate(tc.te, tc.te_data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServeIndex(tc.w, tc.req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if string(responseBody) != ts {
		t.Errorf("Expected data to be equal to ts: %s", string(responseBody))
	    }
	})
    }

    t.Run("redirect", func(t *testing.T){
	if err := testRedirect("/akldfjk", Th.ServeIndex); err != nil {
	    t.Errorf("expected no error and got %v", err)
	}
    })


}

func TestServeBoard(t *testing.T) {
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    createThreads := []sqlc.CreateThreadParams{
	{
	    Title: "This is the first title",
	    Comment: "This is the first comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 1,
	},
	{
	    Title: "This is the second title",
	    Comment: "This is the second comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 1,
	},
	{
	    Title: "This is the third title",
	    Comment: "This is the third comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 1,
	},
	{
	    Title: "This is the first title",
	    Comment: "This is the first comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 2,
	},
	{
	    Title: "This is the second title",
	    Comment: "This is the second comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 2,
	},
	{
	    Title: "This is the third title",
	    Comment: "This is the third comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 2,
	},
	{
	    Title: "This is the first title",
	    Comment: "This is the first comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 3,
	},
	{
	    Title: "This is the second title",
	    Comment: "This is the second comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 3,
	},
	{
	    Title: "This is the third title",
	    Comment: "This is the third comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 3,
	},
    }

    // populating db with threads that are going to be queried

    for _, tt := range createThreads {
	_, err := Th.q.CreateThread(context.Background(), tt)
	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
    }

    testCases := []struct {
	name 	string
	req 	*http.Request
	w	*httptest.ResponseRecorder
	te	*template.Template
    } {
	{
	    name: "tech",
	    req: httptest.NewRequest(http.MethodGet, "/board/tech", nil),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("board"),
	},
	{
	    name: "sports",
	    req: httptest.NewRequest(http.MethodGet, "/board/sports", nil),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("board"),
	},
	{
	    name: "random",
	    req: httptest.NewRequest(http.MethodGet, "/board/random", nil),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("board"),
	},
    }

    for _, tc := range testCases {
	t.Run(tc.name, func(t *testing.T){
	    id := utils.GetBoardID(tc.name)
	    data, err := utils.GetBoardData(Th.q, id, tc.name)
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    ts , err := stringTemplate(tc.te, data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServeBoard(tc.w, tc.req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if string(responseBody) != ts {
		t.Errorf("Expected data to be equal to ts: %v", ts)
	    }

	})
    }

    t.Run("redirect", func(t *testing.T){
	if err := testRedirect("/board/fdjladlfkd", Th.ServeBoard); err != nil {
	    t.Errorf("expected no error and got %v", err)
	}
    })
}



func TestServeThread(t *testing.T) {
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    createThreads := []sqlc.CreateThreadParams{
	{
	    Title: "This is the first title",
	    Comment: "This is the first comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 1,
	},
	{
	    Title: "This is the second title",
	    Comment: "This is the second comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 2,
	},
	{
	    Title: "This is the third title",
	    Comment: "This is the third comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    BoardID: 3,
	},
    }

    createReplies := []sqlc.CreateReplyParams{
	{
	    Comment: "This is the first comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    ThreadID: 1,
	},
	{
	    Comment: "This is the second comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    ThreadID: 2,
	},
	{
	    Comment: "This is the third comment",
	    Date: strconv.Itoa(int(time.Now().Unix())),
	    ThreadID: 3,
	},
    }

    // populating db with threads that are going to be queried

    for _, tt := range createThreads {
	_, err := Th.q.CreateThread(context.Background(), tt)
	if err != nil {
	    t.Errorf("Expected no errors, got %v", err)
	}
    }

    threads, err := Th.q.GetThreads(context.Background(), 3)
    if err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    for i, thread := range threads {
	createReplies[i].ThreadID = thread.ThreadID
    }

    for _, r := range createReplies {
	_, err := Th.q.CreateReply(context.Background(), r)
	if err != nil {
	    t.Errorf("Expected no errors, got %v", err)
	}
    }


    testCases := []struct {
	name 	string
	id	int
	w	*httptest.ResponseRecorder
	te	*template.Template
    } {
	{
	    name: "tech",
	    id: int(threads[2].ThreadID),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("thread"),
	},
	{
	    name: "sports",
	    id: int(threads[2].ThreadID),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("thread"),
	},
	{
	    name: "random",
	    id: int(threads[2].ThreadID),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("thread"),
	},
    }

    for _, tc := range testCases {
	t.Run(tc.name, func(t *testing.T) {
	    req := httptest.NewRequest(http.MethodGet, "/thread/" + strconv.Itoa(tc.id), nil)
	    data, err := utils.GetThreadData(Th.q, int32(tc.id))
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    ts , err := stringTemplate(tc.te, data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServeThread(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if string(responseBody) != ts {
		t.Errorf("expected response to be equal to template string")
	    }
	})
    }
    t.Run("redirect", func(t *testing.T){
	if err := testRedirect("/thread/fdjladlfkd", Th.ServeThread); err != nil {
	    t.Errorf("expected no error and got %v", err)
	}
    })
}

func TestServeKill(t *testing.T) {
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    testCases := []struct{
	name 	string
	id	int
	w	*httptest.ResponseRecorder
	te	*template.Template
    } {
	{
	    name: "killed thread",
	    id: 12345,
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("kill"),
	},
    }

    for _, tc := range testCases {
	t.Run(tc.name, func(t *testing.T){
	    data := models.KillData{
		Thread_id: tc.id,
	    }

	    req := httptest.NewRequest(http.MethodGet, "/kill/" + strconv.Itoa(tc.id), nil)

	    ts , err := stringTemplate(tc.te, data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServeKill(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if string(responseBody) != ts {
		t.Errorf("expected response to be equal to template string")
	    }
	})
    }
}

func TestServeError(t *testing.T) {
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    testCases := []struct{
	name 	string
	status	int
	w	*httptest.ResponseRecorder
	te	*template.Template
    } {
	{
	    name: "not found",
	    status: http.StatusNotFound,
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("error"),
	},
	{
	    name: "internal server error",
	    status: http.StatusInternalServerError,
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("error"),
	},
	{
	    name: "other status",
	    status: http.StatusForbidden,
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("error"),
	},
    }

    for _, tc := range testCases {
	t.Run(tc.name, func(t *testing.T){
	    data := utils.CreateErrorData(tc.status)

	    req := httptest.NewRequest(http.MethodGet, "/error/" + strconv.Itoa(tc.status), nil)

	    ts , err := stringTemplate(tc.te, data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServeError(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if string(responseBody) != ts {
		t.Errorf("expected response to be equal to template string")
	    }
	})
    }
}


func TestServePost(t *testing.T) {
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    getTestCases := []struct{
	name 	string
	board	string
	w	*httptest.ResponseRecorder
	te	*template.Template
	data	models.PostData
    } {
	{
	    name: "get with no errors",
	    board: "tech",
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("post"),
	    data: models.PostData{
		Title: "",
		Comment: "",
		Board: "tech",
		Errors: [2]models.FormError{
		    {Bool: false, Message: "", Field: "title"},
		    {Bool: false, Message: "", Field: "comment"},
		},
	    },
	},
    }

    for _, tc := range getTestCases {
	t.Run(tc.name, func(t *testing.T){
	    req := httptest.NewRequest(http.MethodGet, "/post/" + tc.board, nil)

	    ts , err := stringTemplate(tc.te, tc.data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServePost(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if string(responseBody) != ts {
		t.Errorf("expected response to be equal to template string")
	    }
	})
    }
//--data-raw 'title=a+new+post&comment=this+is+the+comment+for+the+new+post'

    postTestCases := []struct{
	name 	string
	board	string
	resPath	string
	body	io.Reader
	w	*httptest.ResponseRecorder
	te	*template.Template
	data	models.PostData
    } {
	{
	    name: "post with no errors",
	    board: "tech",
	    resPath: "/board/tech",
	    body: bytes.NewReader([]byte("title=a+new+post+with+a+very+interesting+title&comment=this+comment+is+too+good+for+you+to+understand")),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("post"),
	    data: models.PostData{
		Title: "",
		Comment: "",
		Board: "tech",
		Errors: [2]models.FormError{
		    {Bool: false, Message: "", Field: "title"},
		    {Bool: false, Message: "", Field: "comment"},
		},
	    },
	},
	{
	    name: "post with errors",
	    board: "tech",
	    resPath: "/post/tech",
	    body: bytes.NewReader([]byte("title=a+new+post+with+a+very+interesting+title&comment=")),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("post"),
	    data: models.PostData{
		Title: "a new post with a very interesting title",
		Comment: "",
		Board: "tech",
		Errors: [2]models.FormError{
		    {Bool: false, Message: "", Field: "title"},
		    {Bool: true, Message: "This field is required", Field: "comment"},
		},
	    },
	},
    }
    
    for _, tc := range postTestCases {
	t.Run(tc.name, func(t *testing.T){
	    req := httptest.NewRequest(http.MethodPost, "/post/" + tc.board, tc.body)
	    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	    Th.ServePost(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    ts , err := stringTemplate(tc.te, tc.data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    if tc.name == "post with errors" {
		responseBody := tc.w.Body.String() 

		if responseBody != ts {
		    t.Errorf("expected response to be equal to template string %s", responseBody)
		}
	    } else {
		url, err := res.Location()
		if err != nil {
		    t.Errorf("Expected no errors, got %v", err)
		}

		if url.Path != "/board/" + tc.board {
		    t.Errorf("expected path to be /post/" + tc.board + " but got " + url.Path)
		}
	    }

	})
    }
}

func TestServeReply(t *testing.T) {
    if err := start(); err != nil {
	t.Errorf("expected no error, got %v", err)
    }

    threadParams := sqlc.CreateThreadParams{
	Title: "This is the first title",
	Comment: "This is the first comment",
	Date: strconv.Itoa(int(time.Now().Unix())),
	BoardID: 1,
    }

    _, err := Th.q.CreateThread(context.Background(), threadParams)
    if err != nil {
	    t.Errorf("Expected no errors, got %v", err)
    }

    threads, err := Th.q.GetThreads(context.Background(), 1)
    if err != nil {
	    t.Errorf("Expected no errors, got %v", err)
    }

    thread := threads[0]

    getTestCases := []struct{
	name 	string
	id	int
	w	*httptest.ResponseRecorder
	te	*template.Template
	data	models.ReplyData
    } {
	{
	    name: "get with no errors",
	    id: int(thread.ThreadID),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("reply"),
	    data: models.ReplyData{
		Comment: "",
		Thread_id: int(thread.ThreadID),
		Error: models.FormError{
		    Bool: false, 
		    Message: "", 
		    Field: "comment",
		},
	    },
	},
    }

    for _, tc := range getTestCases {
	t.Run(tc.name, func(t *testing.T){
	    req := httptest.NewRequest(http.MethodGet, "/reply/" + strconv.Itoa(tc.id), nil)

	    ts , err := stringTemplate(tc.te, tc.data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    Th.ServeReply(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    responseBody, err := ioutil.ReadAll(res.Body) 
	    if err != nil {
		t.Errorf("Expected no error, got %v", err)
	    }

	    if cleanString(string(responseBody)) != cleanString(ts) {
		t.Errorf("expected response body to be equal to ts (not working, strings are equal) \n%s \n\n ts: \n%s", cleanString(string(responseBody)), cleanString(ts))
	    }
	})
    }

    postTestCases := []struct{
	name 	string
	id	int
	resPath	string
	body	io.Reader
	w	*httptest.ResponseRecorder
	te	*template.Template
	data	models.ReplyData
    } {
	{
	    name: "post with no errors",
	    id: int(thread.ThreadID),
	    resPath: "/thread/" + strconv.Itoa(int(thread.ThreadID)),
	    body: bytes.NewReader([]byte("comment=this+comment+is+too+good+for+you+to+understand")),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("reply"),
	    data: models.ReplyData{
		Comment: "",
		Thread_id: int(thread.ThreadID),
		Error: models.FormError{Bool: false, Message: "", Field: "comment"},
	    },
	},
	{
	    name: "post with errors",
	    id: int(thread.ThreadID),
	    resPath: "/reply/" + strconv.Itoa(int(thread.ThreadID)),
	    body: bytes.NewReader([]byte("comment=")),
	    w: httptest.NewRecorder(), 
	    te: utils.Serve("reply"),
	    data: models.ReplyData{
		Comment: "",
		Thread_id: int(thread.ThreadID),
		Error: models.FormError{Bool: true, Message: "This field is required", Field: "comment"},
	    },
	},
    }
    
    for _, tc := range postTestCases {
	t.Run(tc.name, func(t *testing.T){
	    req := httptest.NewRequest(http.MethodPost, "/reply/" + strconv.Itoa(tc.id), tc.body)
	    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	    Th.ServeReply(tc.w, req)
	    res := tc.w.Result()
	    defer res.Body.Close()

	    ts , err := stringTemplate(tc.te, tc.data)
	    if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	    }

	    if tc.name == "post with errors" {
		responseBody := tc.w.Body.String() 

		if responseBody != ts {
		    t.Errorf("expected response to be equal to template string %s", responseBody)
		}
	    } else {
		url, err := res.Location()
		if err != nil {
		    t.Errorf("Expected no errors, got %v", err)
		}

		if url.Path != "/thread/" + strconv.Itoa(tc.id) {
		    t.Errorf("expected path to be /reply/" + strconv.Itoa(tc.id) + " but got " + url.Path)
		}
	    }

	})
    }
}












