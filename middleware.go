package apicore

import (
	"errors"
	"net/http"
)

type MiddleWare interface {
	// handler可以通过ctx来控制流程
	Before(ctx Context)

	After(ctx Context)
	// 控制处理顺序
	Index() int
}

var middlewareMap = make([]MiddleWare, 0, 10)

// 添加中间件
func AddMiddleware(generator func() MiddleWare) {
	ware := generator()
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

func loadChan(ctx *conn, r *http.Request) {
	var _md MiddleWare
	var err error
	var handler = ctx.GetHandler()
	for _, _md = range middlewareMap {
		_md.Before(ctx)
		if ctx.isBreak() {
			goto END
		}
	}
	// TODO: 加入业务处理前钩子
	// 验证数据合理性
	getRequest(handler, r)
	err = handler.IsValid()
	// 数据无效则跳过处理
	if err != nil {
		ctx.SetRsp(NewClientErrorResponse(errors.New("无效数据")))
		goto AFTER
	}
	// 开始处理
	ctx.GetHandler().Handle(ctx)

AFTER:
	// 后半段中间件调用
	for i := len(middlewareMap) - 1; i > 0; i-- {
		middlewareMap[i].After(ctx)
		if ctx.isBreak() {
			goto END
		}
	}
END:
	return
}
