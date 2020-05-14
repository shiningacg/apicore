package booter

import (
	"context"
	"errors"
	"net/http"
)

////////////////
// 全局Handler管理

var handleMap = make(map[string]func() *handleBase)

/// 添加handler方法

func AddHandler(path string, hdl func() HandleBase) {
	if _, ok := handleMap[path]; ok {
		panic("出现重复的路径：" + path)
	}
	handleMap[path] = func() *handleBase {
		return &handleBase{path: path, handler: hdl()}
	}
}

/// 获取处理函数
func GetRouter() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var res *Response
		var ctx = request.Context()
		if fun, ok := handleMap[request.RequestURI]; ok {
			hdl := fun()
			getRequest(hdl.handler, request)
			err := hdl.IsValid()
			if err != nil {
				res = NewClientErrorResponse(err)
			} else {
				res = getResponse(LoadMiddleware(ctx, request, hdl.Handle))
			}
		} else {
			res = NewClientErrorResponse(errors.New("未知路径"))
		}
		writer.Write(res.Decode())
	}
}

/////////////////

type Handler interface {
	// 路由路径
	Path() string
	HandleBase
}

type HandleBase interface {
	Handle(ctx context.Context, request *http.Request) context.Context
	IsValid() error
}

////////////////////////
// 对handler的基本实现，简单包装

type handleBase struct {
	handler HandleBase
	path    string
}

func (h *handleBase) Handle(ctx context.Context, request *http.Request) context.Context {
	return h.handler.Handle(ctx, request)
}

func (h *handleBase) Path() string {
	return h.path
}

func (h *handleBase) IsValid() error {
	return h.handler.IsValid()
}
