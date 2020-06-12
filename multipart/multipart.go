package multipart

import (
	"apicore"
)

func init() {
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return &Multipart{}
	})
}

type Multipart struct{}

func (m *Multipart) Before(ctx apicore.Context) {
	form, err := ctx.Raw().MultipartReader()
	if err != nil {
		return
	}
	for {
		item, err := form.NextPart()
		if err != nil {
			break
		}
		if item.FileName() != "" {
			ctx.SetValue(item.FormName(), &File{item})
		}
	}
}

func (m *Multipart) After(ctx apicore.Context) {
	return
}

func (m Multipart) Index() int {
	return 10
}
