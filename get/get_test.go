package get

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGet_getArgs(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{name: "tes1", args: args{input: "/login?ph=1"}, want: map[string]string{"ph": "1"}},
		{name: "tes2", args: args{input: "/danmu/v3/?id=15666&max=1000"}, want: map[string]string{"id": "15666", "max": "1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArgs(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanStruct(t *testing.T) {
	type Test1 struct {
		Name string
		ID   int
		Sex  float64 `json:"s"`
	}
	type args struct {
		target interface{}
		args   map[string]string
	}
	tests := []struct {
		name string
		args args
		want *Test1
	}{
		{name: "test1", args: args{target: &Test1{}, args: map[string]string{"name": "1", "id": "19", "s": "1.3"}}, want: &Test1{"1", 19, 1.3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if scanStruct(tt.args.target, tt.args.args); !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("getArgs() = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func Test_scanStruct2(t *testing.T) {
	input := "http://127.0.0.1:3000/danmu/v3/?id=1&max=1000"
	type Input1 struct {
		Id  int `json:"id"`
		Max int `json:"max"`
	}
	t1 := &Input1{}
	scanStruct2(t1, input)
	want1 := &Input1{Id: 1, Max: 1000}
	if !reflect.DeepEqual(t1, want1) {
		t.Errorf("got %v, want %v", t1, want1)
	}
}

// 测试错误的json是不是会影响get到的数据
func Test_Chan(a *testing.T) {
	type Test1 struct {
		Name string
		ID   int
		Sex  float64 `json:"s"`
	}
	var t = &Test1{}
	var args = map[string]string{"name": "1", "id": "19", "s": "1.3"}
	scanStruct(t, args)
	json.Unmarshal([]byte("nil"), t)
	fmt.Println(t)
}

func Test_getArgs(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{name: "test1", args: args{"http://127.0.0.1:3000/danmu/v3/?id=1&max=1000"}, want: map[string]string{"id": "1", "max": "1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArgs(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
