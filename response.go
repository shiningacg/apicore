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
