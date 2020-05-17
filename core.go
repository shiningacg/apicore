package apicore

import (
	"errors"
	"fmt"
	"net/http"
)

func Run(host string) error {
	return http.ListenAndServe(host, &server{})
}

type server struct{}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var res Response
	var ctx = request.Context()
	// 寻找路径，匹配处理方法
	if fun, ok := handleMap[request.RequestURI]; ok {
		ctx = loadChan(ctx, request, fun())
		res = getResponse(ctx)
	} else {
		res = NewClientErrorResponse(errors.New("未知路径"))
	}
	_, err := writer.Write(res.Decode())
	if err != nil {
		fmt.Println(err.Error())
	}
}
