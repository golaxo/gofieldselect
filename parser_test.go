package gofieldselect

import (
	"testing"

	"github.com/golaxo/gofieldselect/internal/lexer"
)

func TestParse_Nested(t *testing.T) {
	t.Parallel()

	input := "id,name,address(street,number),age"
	p := newParser(lexer.New(input))
	node := p.Parse()

	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	identifiers, ok := node.(Identifiers)
	if !ok {
		t.Fatalf("expected nodes to be Identifiers, got %T", node)
	}

	if len(identifiers) != 4 {
		t.Fatalf("expected 4 top-level nodes, got %d", len(identifiers))
	}

	assertIdent(t, identifiers[0], "id")
	assertIdent(t, identifiers[1], "name")

	// address(street,number)
	addr := identifiers[2]
	if addr.Value != "address" {
		t.Fatalf("expected third ident 'address', got %q", addr.Value)
	}

	addressChildren, ok := addr.Child.(Identifiers)
	if !ok {
		t.Fatalf("expected children to be Identifiers, got %T", addr.Child)
	}

	assertIdent(t, addr, "address")
	assertIdent(t, addressChildren[0], "street")
	assertIdent(t, addressChildren[1], "number")
	assertIdent(t, identifiers[3], "age")
}

func TestParse_Whitespace(t *testing.T) {
	t.Parallel()

	input := "  id ,  name , address ( street , number ) , age  "

	p := newParser(lexer.New(input))
	nodes := p.Parse()

	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	identifiers, ok := nodes.(Identifiers)
	if !ok {
		t.Fatalf("expected nodes to be Identifiers, got %T", nodes)
	}

	if len(identifiers) != 4 {
		t.Fatalf("expected 4 top-level nodes, got %d", len(identifiers))
	}
}

func TestParse_JSONKeyChars(t *testing.T) {
	t.Parallel()

	input := "my-name,1,#"
	p := newParser(lexer.New(input))
	nodes := p.Parse()

	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	identifiers, ok := nodes.(Identifiers)
	if !ok {
		t.Fatalf("expected nodes to be Identifiers, got %T", nodes)
	}

	if len(identifiers) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(identifiers))
	}

	assertIdent(t, identifiers[0], "my-name")
	assertIdent(t, identifiers[1], "1")
	assertIdent(t, identifiers[2], "#")
}

func TestParse_Empty(t *testing.T) {
	t.Parallel()

	l := lexer.New("")
	p := newParser(l)
	nodes := p.Parse()

	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	identifiers, ok := nodes.(Identifiers)
	if !ok {
		t.Fatalf("expected nodes to be Identifiers, got %T", nodes)
	}

	if len(identifiers) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(identifiers))
	}
}

func assertIdent(t *testing.T, ident Identifier, want string) {
	t.Helper()

	if ident.Value != want {
		t.Fatalf("expected ident value %q, got %q", want, ident.Value)
	}
}
