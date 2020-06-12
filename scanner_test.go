package apicore

import (
	"context"
	"reflect"
	"testing"
)

func Test_scanStruct(t *testing.T) {
	type Filed1 struct {
		A string
	}
	type Handler1 struct {
		F1 *Filed1 `json:"filed1"`
		B  string
		C  int
	}
	ctx := context.WithValue(context.TODO(), "filed1", &Filed1{A: "hihao"})
	ctx = context.WithValue(ctx, "b", "ahhh")
	ctx = context.WithValue(ctx, "c", "50")

	want := &Handler1{
		F1: &Filed1{"hihao"},
		B:  "ahhh",
		C:  50,
	}
	type args struct {
		target interface{}
		ctx    context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "t1", args: args{target: &Handler1{}, ctx: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//scanStruct(tt.args.target,tt.args.ctx)
			if !reflect.DeepEqual(want, tt.args.target) {
				t.Fatalf("got %v,want %v", tt.args.target, want)
			}
		})
	}
}
