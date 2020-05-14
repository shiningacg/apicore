package example

import (
	"api-template/booter"
	"context"
	"errors"
	"net/http"
)

func init() {
	booter.AddHandler("/login", func() booter.Handler {
		return &Login{}
	})
}

type Login struct {
	UserName string `json:"user_name"`
	UserPWD  string `json:"user_pwd"`
}

func (l *Login) Handle(ctx context.Context, request *http.Request) context.Context {
	if l.UserName == "shlande" && l.UserPWD == "shiningacg" {
		return booter.SetResponse(ctx, booter.NewSuccessResponse(nil))
	}
	return booter.SetResponse(ctx, booter.NewClientErrorResponse(errors.New("账号错误")))
}

func (l *Login) IsValid() error {
	if l.UserPWD == "" || l.UserName == "" {
		return errors.New("无效的输入")
	}
	return nil
}
