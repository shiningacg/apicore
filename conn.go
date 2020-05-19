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
	return &conn{catch, r, p, false}
}

type conn struct {
	// 存放结果
	catch map[string]interface{}
	req   *http.Request
	p     http.ResponseWriter
	write bool
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

func (c *conn) GetCode() int {
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

// 使用完write方法后，其他设置body的方法都会失效
func (c *conn) Write(p []byte) (n int, err error) {
	c.writeHead()
	return c.p.Write(p)
}

// 返回可供写入的内容
func (c *conn) bytes() []byte {
	if c.write {
		return nil
	}
	return c.GetRsp().Decode()
}

// 设置http响应头
func (c *conn) writeHead() {
	c.p.WriteHeader(c.GetCode())
	for head, content := range c.GetHead() {
		c.p.Header().Set(head, content)
	}
}
