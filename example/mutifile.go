package main

import (
	"apicore"
	"apicore/multipart"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func init() {
	apicore.AddHandler(apicore.NewMatcher("/upload"), func() apicore.Handler {
		return &MultiFile{}
	})
}

type MultiFile struct {
	File *multipart.File `json:"file"`
}

func (m *MultiFile) Handle(conn apicore.Conn) {
	fmt.Println(m.File)
	if m.File == nil {
		bt, err := ioutil.ReadFile("PostFile.html")
		if err != nil {
			fmt.Println(err)
			conn.SetRsp(apicore.NewServerErrorResponse(err))
			return
		}
		conn.SetHead("Content-Type", "text/html;charset=utf-8")
		conn.Write(bt)
		return
	}
	file, err := os.Create(m.File.FileName())
	if err != nil {
		conn.SetRsp(apicore.NewServerErrorResponse(err))
	}
	n, err := io.Copy(file, m.File)
	if err != nil {
		conn.SetRsp(apicore.NewServerErrorResponse(err))
	}
	conn.SetRsp(apicore.NewSuccessResponse(n))
}

func (m *MultiFile) IsValid() error {
	return nil
}
