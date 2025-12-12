// Package token contains all the tokens
package token

const (
	// Illegal when an entry isn't allowed.
	Illegal Type = "Illegal"
	// EOF end of filter value.
	EOF Type = "EOF"

	// Ident Field name.
	Ident Type = "Ident"
	// Separator field separator.
	Separator Type = ","

	Lparen Type = "("
	Rparen Type = ")"
)

type (
	Type string

	// Token holds the actual type and its value.
	Token struct {
		// Type of the token.
		Type Type
		// The actual value for the token.
		Literal string
	}
)
