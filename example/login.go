package main

import (
	"apicore"
	"errors"
)

func init() {
	apicore.AddHandler(apicore.NewMatcher("/login"), func() apicore.Handler {
		return &Login{}
	})
}

type Login struct {
	UserName string `json:"user_name"`
	UserPWD  string `json:"user_pwd"`
}

func (l *Login) Handle(ctx apicore.Conn) {
	ctx.SetHead("Access-Control-Allow-Origin", "127.0.0.1:3000")
	if l.UserName == "shlande" && l.UserPWD == "shiningacg" {
		ctx.SetRsp(apicore.NewSuccessResponse(nil))
		return
	}
	ctx.SetRsp(apicore.NewClientErrorResponse(errors.New("账号错误")))
}

func (l *Login) IsValid() error {
	if l.UserPWD == "" || l.UserName == "" {
		return errors.New("无效的输入")
	}
	return nil
}
