package booter

import (
	"context"
	"net/http"
)

type Handler interface {
	Handle(ctx context.Context, request *http.Request) context.Context
	IsValid() error
}

var handleMap = make(map[string]func() Handler)

// 添加handler方法
func AddHandler(path string, hdl func() Handler) {
	if _, ok := handleMap[path]; ok {
		panic("出现重复的路径：" + path)
	}
	handleMap[path] = hdl
}
