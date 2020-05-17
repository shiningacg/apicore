package main

import (
	"apicore"
	"context"
	"errors"
	"net/http"
)

func init() {
	apicore.AddHandler("/login", func() apicore.Handler {
		return &Login{}
	})
}

type Login struct {
	UserName string `json:"user_name"`
	UserPWD  string `json:"user_pwd"`
}

func (l *Login) Handle(ctx context.Context, request *http.Request) context.Context {
	ctx = apicore.SetHead(ctx, "Access-Control-Allow-Origin", "127.0.0.1:3000")
	if l.UserName == "shlande" && l.UserPWD == "shiningacg" {
		return apicore.SetResponse(ctx, apicore.NewSuccessResponse(nil))
	}
	return apicore.SetResponse(ctx, apicore.NewClientErrorResponse(errors.New("账号错误")))
}

func (l *Login) IsValid() error {
	if l.UserPWD == "" || l.UserName == "" {
		return errors.New("无效的输入")
	}
	return nil
}
