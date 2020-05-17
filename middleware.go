package apicore

import (
	"context"
	"errors"
	"net/http"
)

type MiddleWare interface {
	// handler可以通过ctx来控制流程
	Before(ctx context.Context, request *http.Request) context.Context

	After(ctx context.Context, request *http.Request) context.Context
	// 控制处理顺序
	Index() int
}

var middlewareMap = make([]MiddleWare, 0, 10)

// 添加中间件
func AddMiddleware(ware MiddleWare) {
	if len(middlewareMap) == 0 || middlewareMap[len(middlewareMap)-1].Index() <= ware.Index() {
		middlewareMap = append(middlewareMap, ware)
		return
	}
	for i, md := range middlewareMap {
		if md.Index() <= ware.Index() {
			continue
		}
		front := middlewareMap[:i]
		back := append([]MiddleWare{ware}, middlewareMap[i:]...)
		middlewareMap = append(front, back...)
		break
	}
}

// 跳出链式操作
func Break(ctx context.Context) context.Context {
	return context.WithValue(ctx, "SYS_BREAK", true)
}

func isBreak(ctx context.Context) bool {
	return ctx.Value("SYS_BREAK") != nil
}

func loadChan(ctx context.Context, r *http.Request) context.Context {
	var _md MiddleWare
	var err error
	var handler = GetHandler(ctx)
	for _, _md = range middlewareMap {
		ctx = _md.Before(ctx, r)
		if isBreak(ctx) {
			goto END
		}
	}
	// TODO: 加入业务处理前钩子

	// 验证数据合理性
	getRequest(handler, r)
	err = handler.IsValid()
	// 数据无效则跳过处理
	if err != nil {
		ctx = SetResponse(ctx, NewClientErrorResponse(errors.New("无效数据")))
		goto AFTER
	}

	// 开始处理
	ctx = handler.Handle(ctx, r)

AFTER:
	// 后半段中间件调用
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
