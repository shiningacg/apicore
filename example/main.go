package example

import "api-template/booter"

func Main() {
	err := booter.Run(":3000")
	if err != nil {
		panic(err)
	}
}
