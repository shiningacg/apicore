package apicore

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Response interface {
	Decode() []byte
}

type response struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (r *response) Decode() []byte {
	byte, _ := json.Marshal(*r)
	return byte
}

func NewSuccessResponse(data interface{}) Response {
	return &response{
		StatusCode: 200,
		Message:    "ok",
		Data:       data,
	}
}

func NewClientErrorResponse(err error) Response {
	return &response{
		StatusCode: 300,
		Message:    err.Error(),
	}
}

func NewServerErrorResponse(err error) Response {
	return &response{
		StatusCode: 500,
		Message:    "服务器出现错误",
	}
}

func getRequest(target interface{}, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.Decode(target)
}

func SetResponse(ctx context.Context, response Response) context.Context {
	return context.WithValue(ctx, "SYS_RESPONSE", response)
}
func getResponse(ctx context.Context) Response {
	l := ctx.Value("SYS_RESPONSE")
	if l == nil {
		return NewServerErrorResponse(errors.New("no response"))
	}
	return l.(Response)
}

func getHead(ctx context.Context) map[string]string {
	heads := ctx.Value("SYS_HEAD")
	if heads == nil {
		return make(map[string]string, 0)
	}
	return heads.(map[string]string)
}
func SetHead(ctx context.Context, head string, content string) context.Context {
	var heads map[string]string
	if ctx.Value("SYS_HEAD") == nil {
		heads = make(map[string]string)
		heads[head] = content
		return context.WithValue(ctx, "SYS_HEAD", heads)
	}
	heads = ctx.Value("SYS_HEAD").(map[string]string)
	heads[head] = content
	return ctx
}

func getStatusCode(ctx context.Context) int {
	code := ctx.Value("SYS_STATUSCODE")
	if code == nil {
		return 500
	}
	return code.(int)
}
func SetStatusCode(ctx context.Context, code int) context.Context {
	return context.WithValue(ctx, "SYS_STATUSCODE", code)
}
