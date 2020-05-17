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
	var ctx = request.Context()
	// 寻找路径，匹配处理方法
	if fun, ok := handleMap[getPath(request.RequestURI)]; ok {
		ctx = setHandler(ctx, fun())
		ctx = loadChan(ctx, request)
		res = getResponse(ctx)
	} else {
		res = NewClientErrorResponse(errors.New("未知路径"))
	}
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
