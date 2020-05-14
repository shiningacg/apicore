package booter

import "net/http"

func Serve() {

}

type Server struct{}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}
