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
	GetHandler() Handler
	Break()
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

func NewConn(r *http.Request, p http.ResponseWriter) *conn {
	var caches = make(map[string]interface{})
	return &conn{caches, r, p, false}
}

type conn struct {
	// 存放结果
	caches map[string]interface{}
	req    *http.Request
	p      http.ResponseWriter
	write  bool
}

func (c *conn) Value(key string) interface{} {
	return c.caches[key]
}

func (c *conn) SetValue(key string, value interface{}) {
	c.caches[key] = value
}

func (c *conn) Break() {
	c.caches["SYS_BREAK"] = nil
}

func (c *conn) isBreak() bool {
	if _, has := c.caches["SYS_BREAK"]; has {
		return true
	}
	return false
}

func (c *conn) setHandler(handler Handler) {
	c.caches["SYS_HANDLER"] = handler
}

func (c *conn) GetHandler() Handler {
	return c.caches["SYS_HANDLER"].(Handler)
}

func (c *conn) SetCode(code int) {
	c.caches["SYS_CODE"] = code
}

func (c *conn) GetCode() int {
	if code, has := c.caches["SYS_CODE"]; has {
		return code.(int)
	}
	return 200
}

func (c *conn) SetHead(head string, content string) {
	if heads, has := c.caches["SYS_HEAD"]; has {
		heads.(map[string]string)[head] = content
		return
	}
	var heads = make(map[string]string)
	heads[head] = content
	c.SetValue("SYS_HEAD", heads)
}

func (c *conn) GetHead() map[string]string {
	if heads, has := c.caches["SYS_HEAD"]; has {
		return heads.(map[string]string)
	}
	return nil
}

func (c *conn) SetRsp(r Response) {
	c.caches["SYS_RESPONSE"] = r
}

func (c *conn) GetRsp() Response {
	if r, has := c.caches["SYS_RESPONSE"]; has {
		return r.(Response)
	}
	return NewSuccessResponse(nil)
}

func (c *conn) Raw() *http.Request {
	return c.req
}

// 使用完write方法后，其他设置body的方法都会失效
func (c *conn) Write(p []byte) (n int, err error) {
	c.writeHead()
	c.write = true
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
	// 控制只写入一次
	if c.write {
		return
	}
	for head, content := range c.GetHead() {
		c.p.Header().Set(head, content)
	}
	// 清空head，防止重复写入
	c.p.WriteHeader(c.GetCode())
}
