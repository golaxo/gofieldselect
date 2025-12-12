package examples

import (
	"encoding/json"
)

type (
	User struct {
		Name    string   `json:"name,omitempty"`
		Surname string   `json:"surname,omitempty"`
		Age     int      `json:"age,omitempty"`
		Address *Address `json:"address,omitempty"`
	}

	Address struct {
		Street string `json:"street,omitempty"`
		Number int    `json:"number,omitempty"`
	}
)

func DefaultUserExample() User {
	return User{
		Name:    "John",
		Surname: "Doe",
		Age:     18,
		Address: &Address{
			Street: "Main",
			Number: 42,
		},
	}
}

func ToIndentedJson(v any) string {
	res, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(res)
}
