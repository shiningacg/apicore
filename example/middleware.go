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

type Middleware struct{}

const MWN = "ip-watcher"

func GetIP(remoteAddr string) string {
	// ipv4
	if ip := strings.Split(remoteAddr, ":"); len(ip) == 2 {
		return ip[0]
	}
	// ipv6
	return strings.Split(remoteAddr, "]")[0][1:]
	//
}

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

func (m *Middleware) Index() int {
	return 1
}

func Value(ctx context.Context) bool {
	return ctx.Value(MWN).(bool)
}
