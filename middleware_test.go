package apicore

import (
	"testing"
)

func TestAddMiddleware(t *testing.T) {
	input := map[string]int{"t1": 1, "t2": 4, "t3": 3, "t4": 2}
	want := []string{"t1", "t4", "t3"}
	for key, value := range input {
		AddMiddleware(func() MiddleWare {
			return t_middleware{name: key, index: value}
		})
	}
	index := 0
	for i, _ := range want {
		if name := middlewareMap[i].(t_middleware).name; name != want[index] {
			t.Log("want:" + want[index] + ",got:" + name)
		}
		index++
	}
}

type t_middleware struct {
	name  string
	index int
}

func (t t_middleware) Before(ctx Context) {
	if t.index == 3 {
		ctx.Break()
	}
}

func (t t_middleware) After(ctx Context) {
	panic("implement me")
}

func (t t_middleware) Index() int {
	return t.index
}
