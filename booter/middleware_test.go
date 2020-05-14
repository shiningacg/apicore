package booter

import (
	"context"
	"net/http"
	"testing"
)

func TestAddMiddleware(t *testing.T) {
	input := map[string]int{"t1": 1, "t2": 4, "t3": 3, "t4": 2}
	want := []string{"t1", "t4", "t3", "t2"}
	for key, value := range input {
		AddMiddleware(t_middleware{name: key, index: value})
	}
	index := 0
	for _, md := range middlewareMap {
		if name := md.(t_middleware).name; name != want[index] {
			t.Log("want:" + want[index] + ",got:" + name)
		}
		index++
	}
}

type t_middleware struct {
	name  string
	index int
}

func (t t_middleware) Before(ctx context.Context, request *http.Request) context.Context {
	panic("implement me")
}
func (t t_middleware) After(ctx context.Context, request *http.Request) context.Context {
	panic("implement me")
}
func (t t_middleware) Index() int {
	return t.index
}
