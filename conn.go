package apicore

import (
	"io"
	"net/http"
)

// 中间件使用的
type Context interface {
	// 中间件存放结果的地方
	Value(key string) interface{}
	SetValue(key string, value interface{})
	Conn
}

type Conn interface {
	// 基本操作,如果中间件调用则必须break，不然无效
	SetCode(code int)
	SetHead(head string, content string)
	SetRsp(r Response)
	// 返回原始请求
	Raw() *http.Request
	// 写数据，如果调用那么SetRsp会失效
	io.Writer
}

// 流程控制：Before->Handler(writer)->After

func NewConn(r *http.Request, p http.ResponseWriter) Context {
	var catch = make(map[string]interface{})
	return &conn{catch, r, func(bytes []byte) (int, error) {
		return p.Write(bytes)
	}}
}

type conn struct {
	// 存放结果
	catch map[string]interface{}
	req   *http.Request
	r     func([]byte) (int, error)
}

func (c *conn) Value(key string) interface{} {
	return c.catch[key]
}

func (c *conn) SetValue(key string, value interface{}) {
	c.catch[key] = value
}

func (c *conn) SetCode(code int) {
	c.catch["SYS_CODE"] = code
}

func (c *conn) GetCode(code int) int {
	if code, has := c.catch["SYS_CODE"]; has {
		return code.(int)
	}
	return 200
}

func (c *conn) SetHead(head string, content string) {
	if heads, has := c.catch["SYS_HEAD"]; has {
		heads.(map[string]string)[head] = content
		return
	}
	var heads = make(map[string]string)
	heads[head] = content
	c.SetValue("SYS_HEAD", heads)
}

func (c *conn) GetHead() map[string]string {
	if heads, has := c.catch["SYS_HEAD"]; has {
		return heads.(map[string]string)
	}
	return nil
}

func (c *conn) SetRsp(r Response) {
	c.catch["SYS_RESPONSE"] = r
}

func (c *conn) GetRsp() Response {
	if r, has := c.catch["SYS_RESPONSE"]; has {
		return r.(Response)
	}
	return nil
}

func (c *conn) Raw() *http.Request {
	return c.Raw()
}

func (c *conn) Write(p []byte) (n int, err error) {
	return c.r(p)
}
