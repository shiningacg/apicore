package main

import (
	"apicore"
	"context"
	"errors"
	"strings"
)

func init() {
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return &Middleware{}
	})
}

// middleware是单例模式，因此如果要存储值，需要保证线程安全
type Middleware struct{}

// 中间件名称，全局唯一
const MWN = "ip-watcher"

func (m *Middleware) Before(ctx apicore.Context) {
	if GetIP(ctx.Raw().RemoteAddr) == "127.0.0.1" {
		ctx.SetRsp(apicore.NewClientErrorResponse(errors.New("127")))
		ctx.SetValue(MWN, true)
		ctx.Break()
		return
	}
	return
}

func (m *Middleware) After(ctx apicore.Context) {
	return
}

// 优先级
func (m *Middleware) Index() int {
	return 1
}

func Value(ctx context.Context) bool {
	return ctx.Value(MWN).(bool)
}

// 功能实现
func GetIP(remoteAddr string) string {
	// ipv4
	if ip := strings.Split(remoteAddr, ":"); len(ip) == 2 {
		return ip[0]
	}
	// ipv6
	return strings.Split(remoteAddr, "]")[0][1:]
	//
}
