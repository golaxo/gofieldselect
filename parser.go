package gofieldselect

import (
	"slices"

	"codeberg.org/manuelarte/gofieldselect/internal/lexer"
	"codeberg.org/manuelarte/gofieldselect/internal/token"
)

type parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []error
}

// New creates a new Parser based on a Lexer.
func newParser(l *lexer.Lexer) *parser {
	p := &parser{
		l:      l,
		errors: make([]error, 0),
	}
	// Initialize tokens
	p.nextToken()
	p.nextToken()

	return p
}

// Errors return any parsing errors encountered.
func (p *parser) Errors() []error {
	return slices.Clone(p.errors)
}

// Parse parses the input stream into a list of Nodes (top-level fields).
func (p *parser) Parse() Node {
	return p.parseFields()
}

func (p *parser) nextToken() { p.curToken = p.peekToken; p.peekToken = p.l.NextToken() }

// parseFields parses a comma-separated list of fields until a right parenthesis or EOF.
// It assumes p.curToken is positioned at the first token of the list (which can be Ident, Rparen, or EOF).
func (p *parser) parseFields() Node {
	nodes := make([]Identifier, 0)

	for p.curToken.Type != token.EOF && p.curToken.Type != token.Rparen {
		if p.curToken.Type != token.Ident {
			// unexpected token; attempt to recover by skipping until next separator, rparen, or EOF
			p.errors = append(p.errors, ErrExpectedIdentifier)
			p.synchronize()

			if p.curToken.Type == token.Separator {
				p.nextToken()
			}

			continue
		}

		n := p.parseField()
		nodes = append(nodes, n)

		// After a field, the current token is expected to be either ',' or ')' or EOF
		switch p.curToken.Type {
		case token.Separator:
			// consume ',' to move to the next field (which should be Ident or end)
			p.nextToken()
		case token.Rparen, token.EOF:
			// list ends; loop condition will break
		case token.Ident, token.Illegal, token.Lparen:
			// missing comma between identifiers; record error but continue without consuming to avoid infinite loop
			p.errors = append(p.errors, ErrMissingSeparatorBetweenIdentifiers)
		default:
			// Any other token, just advance to keep progress
			p.nextToken()
		}
	}

	return Identifiers(nodes)
}

// parseField parses a single identifier with optional nested children in parentheses.
// Precondition: p.curToken.Type == token.Ident
// Postcondition: p.curToken will be the token following the field (','/')'/EOF).
func (p *parser) parseField() Identifier {
	ident := Identifier{Value: p.curToken.Literal, Child: AllIdentifiers{}}

	if p.peekToken.Type == token.Lparen {
		// consume '(' and move inside
		p.nextToken() // move to '('
		p.nextToken() // move to first token inside children

		// parse children until we hit ')'
		child := p.parseFields()
		ident.Child = child

		if p.curToken.Type == token.Rparen {
			// consume ')'
			p.nextToken()
		} else {
			p.errors = append(p.errors, ErrExpectedClosingParenthesis)
		}
	} else {
		// move past the identifier; caller will handle separators
		p.nextToken()
	}

	return ident
}

// synchronize advances tokens until a safe point (comma, right parenthesis, or EOF).
func (p *parser) synchronize() {
	for p.curToken.Type != token.Separator && p.curToken.Type != token.Rparen && p.curToken.Type != token.EOF {
		p.nextToken()
	}
}
