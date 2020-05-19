package main

/*
	这个文件是用来测试write是否按照正常的流程运行的
*/

import "apicore"

func init() {
	apicore.AddHandler(apicore.NewMatcher("/stream", "GET"), func() apicore.Handler {
		return &Stream{}
	})
}

type Stream struct {
}

func (s *Stream) Handle(conn apicore.Conn) {
	conn.SetCode(300)
	conn.SetHead("X-Hello", "Ok")
	conn.Write([]byte("hellow"))
}

func (s *Stream) IsValid() error {
	return nil
}
