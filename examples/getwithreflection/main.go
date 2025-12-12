package main

import (
	"fmt"

	"examples"

	"codeberg.org/manuelarte/gofieldselect"
)

func main() {
	fieldSelection := "name,surname"
	src := examples.DefaultUserExample()

	n, err := gofieldselect.Parse(fieldSelection)
	if err != nil {
		panic(err)
	}
	selected, err := gofieldselect.GetWithReflection(n, src)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Source: \n %+v\n", examples.ToIndentedJson(src))
	fmt.Printf("Selected: \n%+v\n", examples.ToIndentedJson(selected))
}
