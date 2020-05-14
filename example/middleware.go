package example

import (
	"api-template/booter"
	"context"
	"errors"
	"net/http"
	"strings"
)

func init() {
	booter.AddMiddleware(&Middleware{})
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

func (m *Middleware) Before(ctx context.Context, request *http.Request) context.Context {
	if GetIP(request.RemoteAddr) == "127.0.0.1" {
		ctx = booter.SetResponse(ctx, booter.NewClientErrorResponse(errors.New("127")))
		ctx = context.WithValue(ctx, MWN, true)
		return booter.Break(ctx)
	}
	return ctx
}

func (m *Middleware) After(ctx context.Context, request *http.Request) context.Context {
	return ctx
}

func (m *Middleware) Index() int {
	return 1
}

func Value(ctx context.Context) bool {
	return ctx.Value(MWN).(bool)
}
