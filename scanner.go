package apicore

import (
	"reflect"
	"strconv"
	"strings"
)

func init() {
	AddMiddleware(func() MiddleWare {
		return Scanner{}
	})
}

type Scanner struct{}

func (s Scanner) Before(ctx Context) {
	scanStruct(ctx.GetHandler(), ctx)
}

func (s Scanner) After(ctx Context) {
	return
}

func (s Scanner) Index() int {
	return 10000
}

func scanStruct(target interface{}, ctx Context) {
	v := reflect.ValueOf(target).Elem()
	t := reflect.TypeOf(target).Elem()
	count := v.NumField()
	for i := 0; i < count; i++ {
		// 判断filed的类型是不是指针
		vl := v.Field(i)
		tp := t.Field(i)
		var value interface{}
		// 字段名称
		Name := tp.Name
		// Tag名
		Tag := tp.Tag.Get("json")
		// 寻找字段
		if temp := ctx.Value(strings.ToLower(Name)); temp != nil {
			value = temp
		} else if temp := ctx.Value(Tag); temp != nil {
			value = temp
		} else {
			continue
		}
		val, ok := value.(string)
		// 如果不成功，那么为特殊类型，直接通过指针进行操作
		if !ok {
			_v := reflect.ValueOf(value)
			if _v.Kind() == reflect.Ptr && _v.Type() == tp.Type {
				vl.Set(_v)
			}
			continue
		}
		// 如果是字符串类型，那么进行可能的类型转化
		switch vl.Kind() {
		case reflect.Int:
			temp, err := strconv.ParseInt(val, 0, 64)
			if err != nil {
				break
			}
			v.Field(i).SetInt(temp)
		case reflect.String:
			v.Field(i).SetString(val)
		case reflect.Float64:
			temp, err := strconv.ParseFloat(val, 64)
			if err != nil {
				break
			}
			v.Field(i).SetFloat(temp)
		}
	}
}
