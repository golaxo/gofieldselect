package lexer

import (
	"testing"

	"codeberg.org/manuelarte/gofieldselect/internal/token"
)

func TestNextToken_Nested(t *testing.T) {
	t.Parallel()

	input := "id,name,address(street,number),age"

	expected := []token.Token{
		{Type: token.Ident, Literal: "id"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "name"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "address"},
		{Type: token.Lparen, Literal: "("},
		{Type: token.Ident, Literal: "street"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "number"},
		{Type: token.Rparen, Literal: ")"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "age"},
		{Type: token.EOF, Literal: ""},
	}

	l := New(input)

	for i, tt := range expected {
		tok := l.NextToken()

		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.Type, tok.Type)
		}

		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.Literal, tok.Literal)
		}
	}
}

func TestNextToken_Whitespace(t *testing.T) {
	t.Parallel()

	input := "  id ,  name , address ( street , number ) , age  "

	expected := []token.Token{
		{Type: token.Ident, Literal: "id"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "name"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "address"},
		{Type: token.Lparen, Literal: "("},
		{Type: token.Ident, Literal: "street"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "number"},
		{Type: token.Rparen, Literal: ")"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "age"},
		{Type: token.EOF, Literal: ""},
	}

	l := New(input)

	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.Type || tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - expected (%q,%q), got (%q,%q)", i, tt.Type, tt.Literal, tok.Type, tok.Literal)
		}
	}
}

func TestNextToken_Empty(t *testing.T) {
	t.Parallel()

	l := New("")

	tok := l.NextToken()
	if tok.Type != token.EOF {
		t.Fatalf("expected EOF, got %q (%q)", tok.Type, tok.Literal)
	}
}

func TestNextToken_JSONKeyChars(t *testing.T) {
	t.Parallel()

	input := "my-name,1,#"

	expected := []token.Token{
		{Type: token.Ident, Literal: "my-name"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "1"},
		{Type: token.Separator, Literal: ","},
		{Type: token.Ident, Literal: "#"},
		{Type: token.EOF, Literal: ""},
	}

	l := New(input)

	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.Type || tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - expected (%q,%q), got (%q,%q)", i, tt.Type, tt.Literal, tok.Type, tok.Literal)
		}
	}
}

func TestNextToken_Tab(t *testing.T) {
	t.Parallel()

	expected := []token.Token{
		{Type: token.Illegal, Literal: "\t"},
		{Type: token.EOF, Literal: ""},
	}

	l := New("\t")

	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.Type || tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - expected (%q,%q), got (%q,%q)", i, tt.Type, tt.Literal, tok.Type, tok.Literal)
		}
	}
}
