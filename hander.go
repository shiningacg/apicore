package apicore

type Handler interface {
	Handle(conn Conn)
	IsValid() error
}

var handleMap = make(map[string]func() Handler)

// 添加handler方法
func AddHandler(path string, hdl func() Handler) {
	if _, ok := handleMap[path]; ok {
		panic("出现重复的路径：" + path)
	}
	handleMap[path] = hdl
}
