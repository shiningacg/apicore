package get

import (
	"github.com/shiningacg/apicore"
	"net/url"
)

func init() {
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return &Get{}
	})
}

type Get struct{}

func (g *Get) Before(ctx apicore.Context) {
	if ctx.Raw().Method == "GET" {
		u, err := url.Parse(ctx.Raw().RequestURI)
		if err != nil {
			return
		}
		m := u.Query()
		for name, value := range m {
			ctx.SetValue(name, value[0])
		}
	}
}

func (g *Get) After(ctx apicore.Context) {
	return
}

func (g *Get) Index() int {
	return -1
}
