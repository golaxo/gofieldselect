// Package main example on how to use gofieldselect.GetValue to transform a DAO object to a DTO with field selection.
package main

import (
	"fmt"

	"examples"

	"github.com/golaxo/gofieldselect"
)

type (
	userDao struct {
		Name    string     `json:"name"`
		Surname string     `json:"surname"`
		Age     int        `json:"age"`
		Address addressDao `json:"address"`
	}

	addressDao struct {
		Street string `json:"street"`
		Number int    `json:"number"`
	}
)

func main() {
	user := userDao{
		Name:    "John",
		Surname: "Doe",
		Age:     18,
		Address: addressDao{
			Street: "Example",
			Number: 1,
		},
	}

	fieldSelection, err := gofieldselect.Parse("name,address(street)")
	if err != nil {
		panic(err)
	}

	var address *examples.Address
	addressFieldSelectionNode, ok := fieldSelection.SelectField("address")
	if ok {
		address = &examples.Address{
			Street: gofieldselect.Get(addressFieldSelectionNode.Child, "street", user.Address.Street),
			Number: gofieldselect.Get(addressFieldSelectionNode.Child, "number", user.Address.Number),
		}
	}

	dto := examples.User{
		Name:    gofieldselect.Get(fieldSelection, "name", user.Name),
		Surname: gofieldselect.Get(fieldSelection, "surname", user.Surname),
		Age:     gofieldselect.Get(fieldSelection, "age", user.Age),
		Address: address,
	}

	fmt.Printf("Source: \n %+v\n", examples.ToIndentedJson(user))
	fmt.Printf("Selected: \n%+v\n", examples.ToIndentedJson(dto))
}
