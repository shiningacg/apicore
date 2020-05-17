package main

import (
	"api-template"
)

func main() {
	err := apicore.Run(":3000")
	if err != nil {
		panic(err)
	}
}
