// Package lexer contains the lexer for gofieldselect
package lexer

import (
	"github.com/golaxo/gofieldselect/internal/token"
)

// Lexer to parse the input.
type Lexer struct {
	input string
	// Current position in input (points to current char)
	position int
	// Current reading position in input (after current char)
	readPosition int
	// current character under examination
	ch byte
}

// New creates a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// Initialize positions so that getChar works correctly
	l.position = 0
	l.readPosition = 0
	l.readChar()

	return l
}

// NextToken returns the next token from the input stream.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case ',':
		tok = newToken(token.Separator, l.ch)
	case '(':
		tok = newToken(token.Lparen, l.ch)
	case ')':
		tok = newToken(token.Rparen, l.ch)
	case '\n', '\t', '\r':
		tok = newToken(token.Illegal, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		// Ident: any JSON key characters until delimiter or whitespace
		if isIdentChar(l.ch) {
			literal := l.readIdentifier()
			tok.Type = token.Ident
			tok.Literal = literal

			return tok
		}

		tok = newToken(token.Illegal, l.ch)
	}

	l.readChar()

	return tok
}

func newToken(t token.Type, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isIdentChar(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

// isWhitespace reports whether the given byte is a whitespace we should skip between tokens.
func isWhitespace(ch byte) bool {
	return ch == ' '
}

// isDelimiter is any character that separates tokens and is not part of an identifier.
func isDelimiter(ch byte) bool {
	return ch == ',' || ch == '(' || ch == ')' || ch == 0
}

// isIdentChar reports whether the byte can be part of an identifier (JSON key) in this grammar.
// We allow any non-delimiter, non-whitespace character sequence.
func isIdentChar(ch byte) bool {
	return !isDelimiter(ch) && !isWhitespace(ch)
}
