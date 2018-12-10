package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asdine/storm"
)

const (
	OperationSuccess = "success"
	OperationFailed  = "failed"
	OperationGET     = "get"
	OperationDELETE  = "delete"
	OperationUPDATE  = "update"
	OperationCREATE  = "create"
)

type Operation struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

// APIMethods ...
type APIMethods interface {
	GET(ctx *Context)
	PUT(ctx *Context)
	POST(ctx *Context)
	DELETE(ctx *Context)
}

type BaseView interface {
	APIMethods
	Dispatch(string, *storm.DB) http.HandlerFunc
}

type View struct {
	BaseView
}

type Response struct {
	Op   *Operation  `json:"op"`
	Data interface{} `json:"data"`
}
type Context struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	Operation *Operation
	Response  *Response
	DB        *storm.DB
	Token     string
}

func (ctx *Context) InitBucket(bucketName string) {
	err := ctx.DB.Init(bucketName)
	fmt.Println(err)
}

// Dispatch merges handlers into one
func (API *View) Dispatch(bucketName string, db *storm.DB) http.HandlerFunc {
	db.Init(bucketName)
	return func(w http.ResponseWriter, r *http.Request) {
		handleRequest(API, w, r, db)
	}
}

// NewResponse constructor
func NewResponse(op *Operation, data interface{}) *Response {
	res := &Response{Op: op}
	if op.Status == OperationFailed {
		return res
	}
	res.Data = data
	return res
}

// JSONResponse ...
func JSONResponse(ctx *Context) {
	jsonVar, err := json.Marshal(ctx.Response)

	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Writer.Header().Set("content-type", "application/json")
	ctx.Writer.Write(jsonVar)

}

// NewView constructor
func NewView(v BaseView) *View {
	return &View{BaseView: v}
}

// NewOperation iniaties an Operation instance with defaults
func NewOperation(name string) *Operation {
	return &Operation{Name: name, Status: OperationSuccess, Error: ""}
}

func (op *Operation) Fail(er error) {
	op.Status = OperationFailed
	op.Error = er.Error()
}

func handleRequest(source APIMethods, w http.ResponseWriter, r *http.Request, db *storm.DB) {
	ctx := &Context{
		Writer:  w,
		Request: r,
		DB:      db,
	}

	if r.Method == "GET" {
		ctx.Operation = NewOperation(OperationGET)
		source.GET(ctx)
	} else if r.Method == "POST" {
		ctx.Operation = NewOperation(OperationCREATE)
		source.POST(ctx)
	} else if r.Method == "PUT" {
		ctx.Operation = NewOperation(OperationUPDATE)
		source.PUT(ctx)
	} else if r.Method == "DELETE" {
		ctx.Operation = NewOperation(OperationDELETE)
		source.DELETE(ctx)
	}

	JSONResponse(ctx)
}
