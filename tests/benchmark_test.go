package tests

import (
	"testing"

	"github.com/golaxo/gofieldselect"
)

type user struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

type user2 struct {
	Name    *string `json:"name,omitempty"`
	Surname *string `json:"surname,omitempty"`
	Age     *int    `json:"age,omitempty"`
}

//nolint:gochecknoglobals // https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
var (
	other  user
	other2 user2
)

func Benchmark_GetWithReflection(b *testing.B) {
	u := user{
		Name:    "John",
		Surname: "Doe",
		Age:     20,
	}

	fields, err := gofieldselect.Parse("name")
	if err != nil {
		b.Fatal(err)
	}

	for range b.N {
		other, _ = gofieldselect.GetWithReflection(fields, u)
	}
}

func Benchmark_Get(b *testing.B) {
	u := user{
		Name:    "John",
		Surname: "Doe",
		Age:     20,
	}

	fields, err := gofieldselect.Parse("name")
	if err != nil {
		b.Fatal(err)
	}

	for range b.N {
		other2 = user2{
			Name:    gofieldselect.Get(fields, "name", &u.Name),
			Surname: gofieldselect.Get(fields, "name", &u.Surname),
			Age:     gofieldselect.Get(fields, "name", &u.Age),
		}
	}
}
