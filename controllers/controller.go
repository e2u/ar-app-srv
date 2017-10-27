package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"

	"e2u.io/ar-app-srv/app"
	"github.com/e2u/goboot"
	"github.com/go-chi/chi/middleware"
)

const (
	RespCodeTag = "resp_code"
	ErrorsTag   = "errors"
	RemarkTag   = "remark"

	ResponseSuccess = 1000
	ResponseError   = 9000

	RemarkSuccess = "success"
	RemarkError   = "error"
)

type Controller struct {
	AppContext *app.AppContext
	*Response
	params *goboot.Params
}

type Response struct {
	errors     []interface{}
	outPutResp *sync.Map
}

func NewResponse() *Response {
	return &Response{
		outPutResp: &sync.Map{},
	}
}

func (resp *Response) Error(e interface{}) {
	resp.PutResp(RespCodeTag, ResponseError)
	resp.PutResp(RemarkTag, RemarkError)
	resp.errors = append(resp.errors, e)
}

func (resp *Response) Success() {
	resp.PutResp(RespCodeTag, ResponseSuccess)
	resp.PutResp(RemarkTag, RemarkSuccess)
}

func (resp *Response) PutResp(key string, value interface{}) {
	resp.outPutResp.Store(key, value)
}

func (resp *Response) RenderJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if _, ok := resp.outPutResp.Load(RespCodeTag); !ok {
		resp.PutResp(RespCodeTag, ResponseSuccess)
	}

	if len(resp.errors) > 0 {
		resp.PutResp(RespCodeTag, ResponseError)
		resp.PutResp(ErrorsTag, resp.errors)
	}
	outMap := make(map[string]interface{})
	resp.outPutResp.Range(func(key interface{}, value interface{}) bool {
		outMap[key.(string)] = value
		return true
	})

	if err := enc.Encode(outMap); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	out := buf.Bytes()
	//goboot.Log.Infof("最终输出内容: [%s] %s", middleware.GetReqID(r.Context()), string(out))

	w.Write(out)
}

func (c *Controller) InitWithMiddle(w http.ResponseWriter, r *http.Request) bool {
	if rd, err := httputil.DumpRequest(r, true); err == nil {
		goboot.Log.Infof("[%s] %s", middleware.GetReqID(r.Context()), rd)
	}

	c.Response = NewResponse()
	c.params = &goboot.Params{}
	goboot.ParseParams(c.params, r)

	return false
}
