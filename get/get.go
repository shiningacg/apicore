package get

import (
	"apicore"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return &Get{}
	})
}

type Get struct{}

func (g *Get) Before(ctx apicore.Context) {
	if ctx.Raw().Method == "GET" {
		handler := ctx.GetHandler()
		scanStruct(handler, ctx.Raw().RequestURI)
	}
}

func (g *Get) After(ctx apicore.Context) {
	return
}

func (g *Get) Index() int {
	return -1
}

func scanStruct(target interface{}, input string) {
	u, err := url.Parse(input)
	if err != nil {
		return
	}
	m := u.Query()
	fmt.Println(m.Get("id"), m.Get("max"))
	v := reflect.ValueOf(target).Elem()
	t := reflect.TypeOf(target).Elem()
	count := v.NumField()
	for i := 0; i < count; i++ {
		var value string
		fieldType := v.Field(i).Type()
		fieldName := t.Field(i).Name
		fieldTag := t.Field(i).Tag.Get("json")
		if temp := m.Get(strings.ToLower(fieldName)); temp != "" {
			value = temp
		}
		if temp := m.Get(fieldTag); temp != "" {
			value = temp
		}
		if value == "" {
			continue
		}
		switch fieldType.Kind() {
		case reflect.Int:
			temp, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				break
			}
			v.Field(i).SetInt(temp)
		case reflect.String:
			v.Field(i).SetString(value)
		case reflect.Float64:
			temp, err := strconv.ParseFloat(value, 64)
			if err != nil {
				break
			}
			v.Field(i).SetFloat(temp)
		}
	}
}
