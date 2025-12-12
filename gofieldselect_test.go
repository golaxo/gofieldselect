package gofieldselect

import (
	"errors"
	"testing"

	"codeberg.org/manuelarte/gofieldselect/internal/lexer"
)

type Address struct {
	Street string `json:"street"`
	Number int    `json:"number"`
}

type User struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type UserPtr struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Age     int      `json:"age"`
	Address *Address `json:"address"`
}

func TestApplyFromNode_FlatFields(t *testing.T) {
	t.Parallel()

	nodes := parse(t, "name,surname")

	src := User{
		Name:    "John",
		Surname: "Doe",
		Age:     18,
		Address: Address{Street: "Main", Number: 42},
	}

	expected := User{
		Name:    src.Name,
		Surname: src.Surname,
	}

	got, err := GetWithReflection(nodes, src)
	if err != nil {
		t.Fatalf("WithReflection returned error: %v", err)
	}

	if got != expected {
		t.Fatalf("expected %+v; got %+v", expected, got)
	}
}

func TestApplyFromNode_PtrField(t *testing.T) {
	t.Parallel()

	nodes := parse(t, "name,surname")

	src := UserPtr{
		Name:    "John",
		Surname: "Doe",
		Age:     18,
		Address: &Address{Street: "Main", Number: 42},
	}

	expected := UserPtr{
		Name:    src.Name,
		Surname: src.Surname,
	}

	got, err := GetWithReflection(nodes, src)
	if err != nil {
		t.Fatalf("WithReflection returned error: %v", err)
	}

	if got != expected {
		t.Fatalf("expected %+v; got %+v", expected, got)
	}
}

func TestApplyFromNode_Nested(t *testing.T) {
	t.Parallel()

	nodes := parse(t, "address(street)")

	src := User{
		Name:    "John",
		Surname: "Doe",
		Age:     18,
		Address: Address{Street: "Example street", Number: 7},
	}

	expected := User{
		Address: Address{Street: "Example street"},
	}

	got, err := GetWithReflection(nodes, src)
	if err != nil {
		t.Fatalf("WithReflection returned error: %v", err)
	}

	if got != expected {
		t.Fatalf("expected %+v; got %+v", expected, got)
	}
}

func TestApplyFromNode_FieldNameWhenNoJSONTag(t *testing.T) {
	t.Parallel()

	type S struct {
		Foo string
		Bar int
	}

	nodes := parse(t, "Foo")
	src := S{Foo: "x", Bar: 3}
	expected := S{Foo: "x"}

	got, err := GetWithReflection(nodes, src)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if got != expected {
		t.Fatalf("expected %+v; got %+v", expected, got)
	}
}

func TestApplyFromNode_PointerInput(t *testing.T) {
	t.Parallel()

	nodes := parse(t, "name")
	src := &User{Name: "Alice", Age: 30}
	expected := &User{Name: "Alice"}

	got, err := GetWithReflection(nodes, src)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if *got != *expected {
		t.Fatalf("expected %+v; got %+v", expected, got)
	}
}

func TestApplyFromNode_NonStructError(t *testing.T) {
	t.Parallel()

	nodes := parse(t, "name")
	_, err := GetWithReflection(nodes, 5)

	expectedErr := TypeNotValidError{}
	if !errors.As(err, &expectedErr) {
		t.Fatalf("expected error for non-struct input")
	}

	if expectedErr.kind.String() != "int" {
		t.Fatalf("expected kind to be int; got %q", expectedErr.kind)
	}
}

// Helper to parse a selection string into nodes.
func parse(t *testing.T, sel string) Node {
	t.Helper()

	p := newParser(lexer.New(sel))

	node := p.Parse()
	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	return node
}
