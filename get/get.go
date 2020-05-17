package get

import (
	"apicore"
	"context"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	apicore.AddMiddleware(&Get{})
}

type Get struct{}

func (g *Get) Before(ctx context.Context, request *http.Request) context.Context {
	if request.Method == "GET" {
		handler := apicore.GetHandler(ctx)
		args := g.getArgs(request.RequestURI)
		scanStruct(handler, args)
	}
	return ctx
}

func (g *Get) After(ctx context.Context, request *http.Request) context.Context {
	return nil
}

func (g *Get) Index() int {
	return -1
}

func (g *Get) getArgs(input string) map[string]string {
	var dict = make(map[string]string)
	rawarg := strings.Split(input, "?")[1]
	args := strings.Split(rawarg, "&")
	for _, arg := range args {
		raw := strings.Split(arg, "=")
		dict[raw[0]] = raw[1]
	}
	return dict
}

func scanStruct(target interface{}, args map[string]string) {
	v := reflect.ValueOf(target).Elem()
	t := reflect.TypeOf(target).Elem()
	count := v.NumField()
	for i := 0; i < count; i++ {
		var value string
		fieldType := v.Field(i).Type()
		fieldName := t.Field(i).Name
		fieldTag := t.Field(i).Tag.Get("json")
		if temp, hasField := args[strings.ToLower(fieldName)]; hasField {
			value = temp
		}
		if temp, hasTag := args[fieldTag]; hasTag {
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
