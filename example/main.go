package example

import (
	"api-template/booter"
	"net/http"
)

func Main() {
	http.ListenAndServe(":3000", &server{})
}

type server struct{}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	booter.GetRouter()(writer, request)
}
