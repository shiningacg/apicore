package apicore

type Handler interface {
	Handle(conn Conn)
	IsValid() error
}

var handleMap = make(map[Matcher]func() Handler)

// 添加handler方法
func AddHandler(matcher Matcher, hdl func() Handler) {
	for m, _ := range handleMap {
		if matcher == m {
			panic("重复添加方法!")
		}
	}
	handleMap[matcher] = hdl
}
