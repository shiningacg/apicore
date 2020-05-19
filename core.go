package apicore

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

func Run(host string) error {
	return http.ListenAndServe(host, &server{})
}

type server struct{}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := NewConn(request, writer)
	// 寻找路径，匹配处理方法
	if fun, ok := handleMap[getPath(request.RequestURI)]; ok {
		ctx.setHandler(fun())
		loadChan(ctx, request)
	} else {
		// TODO:把头部的设置合并到快捷回复中
		ctx.SetCode(300)
		ctx.SetRsp(NewClientErrorResponse(errors.New("未知路径")))
	}
	// 设置头部信息
	ctx.writeHead()
	// 写回复
	writer.Write(ctx.bytes())
}

func getPath(requestURL string) string {
	return strings.Split(requestURL, "?")[0]
}

func setHandler(ctx context.Context, handler Handler) context.Context {
	return context.WithValue(ctx, "SYS_HANDLER", handler)
}
func GetHandler(ctx context.Context) Handler {
	return ctx.Value("SYS_HANDLER").(Handler)
}
