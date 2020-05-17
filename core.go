package apicore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func Run(host string) error {
	return http.ListenAndServe(host, &server{})
}

type server struct{}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var res Response
	var heads map[string]string
	var code int
	var ctx = request.Context()
	// 寻找路径，匹配处理方法
	if fun, ok := handleMap[getPath(request.RequestURI)]; ok {
		ctx = setHandler(ctx, fun())
		ctx = loadChan(ctx, request)
		res = getResponse(ctx)
		code = getStatusCode(ctx)
		heads = getHead(ctx)
	} else {
		// TODO:把头部的设置合并到快捷回复中
		code = 300
		heads = make(map[string]string, 0)
		res = NewClientErrorResponse(errors.New("未知路径"))
	}
	// 设置头部
	writer.WriteHeader(code)
	for hd, ct := range heads {
		writer.Header().Set(hd, ct)
	}
	// 写回复
	_, err := writer.Write(res.Decode())
	if err != nil {
		fmt.Println(err.Error())
	}
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
