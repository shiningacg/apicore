package get

import (
	"apicore"
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
		args := getArgs(ctx.Raw().RequestURI)
		scanStruct(handler, args)
	}
}

func (g *Get) After(ctx apicore.Context) {
	return
}

func (g *Get) Index() int {
	return -1
}

func getArgs(input string) map[string]string {
	var (
		args    []string
		rawargs string
	)
	var dict = make(map[string]string)
	// 如果有参数
	if args = strings.Split(input, "?"); len(rawargs) == 1 {
		return nil
	}
	args = append([]string{}, args[1:]...)
	rawargs = strings.Join(args, "")
	// TODO:解决参数里面有&的情况
	args = strings.Split(rawargs, "&")
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
