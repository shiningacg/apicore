package main

import (
	"github.com/shiningacg/apicore"
	_ "github.com/shiningacg/apicore/get"
)

func main() {
	err := apicore.Run(":3000")
	if err != nil {
		panic(err)
	}
}
