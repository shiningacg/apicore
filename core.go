package apicore

import (
	"errors"
	"net/http"
)

func Run(host string) error {
	return http.ListenAndServe(host, &server{})
}

type server struct{}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := NewConn(request, writer)
	// 寻找路径，匹配处理方法
	for matcher, generator := range handleMap {
		if matcher.Match(request.RequestURI, request.Method) {
			ctx.setHandler(generator())
			loadChan(ctx, request)
			goto END
		}
	}
	// TODO:把头部的设置合并到快捷回复中
	ctx.SetCode(300)
	ctx.SetRsp(NewClientErrorResponse(errors.New("未知路径")))
END:
	// 设置头部信息
	ctx.writeHead()
	// 写回复
	writer.Write(ctx.bytes())
}
