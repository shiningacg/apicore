package booter

import (
	"context"
	"net/http"
)

var middlewareMap = make([]MiddleWare, 0, 10)

func AddMiddleware(ware MiddleWare) {
	middlewareMap = append(middlewareMap, ware)
}

type MiddleWare interface {
	// handler可以通过ctx来控制流程
	Before(ctx context.Context, request *http.Request) context.Context

	After(ctx context.Context, request *http.Request) context.Context
	// 控制处理顺序
	Index() int
}

func Break(ctx context.Context) context.Context {
	return context.WithValue(ctx, "SYS_BREAK", true)
}

func isBreak(ctx context.Context) bool {
	return ctx.Value("SYS_BREAK") != nil
}

func LoadMiddleware(ctx context.Context, r *http.Request, handler func(context.Context, *http.Request) context.Context) context.Context {
	var _md MiddleWare
	for _, _md = range middlewareMap {
		ctx = _md.Before(ctx, r)
		if isBreak(ctx) {
			goto END
		}
	}
	ctx = handler(ctx, r)
	for i := len(middlewareMap) - 1; i > 0; i-- {
		_md = middlewareMap[i]
		ctx = _md.After(ctx, r)
		if isBreak(ctx) {
			goto END
		}
	}
END:
	return ctx
}
