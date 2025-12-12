package gofieldselect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	_ error = new(ParsingError)

	ErrExpectedIdentifier                 = errors.New("expected identifier")
	ErrMissingSeparatorBetweenIdentifiers = errors.New("missing separator between identifiers")
	ErrExpectedClosingParenthesis         = errors.New("expected closing parenthesis")
)

type (
	ParsingError struct {
		errSlice []error
	}

	TypeNotValidError struct {
		kind reflect.Kind
	}
)

func NewParsingError(errSlice []error) ParsingError {
	return ParsingError{errSlice: errSlice}
}

func (pe ParsingError) Error() string {
	ss := make([]string, len(pe.errSlice))
	for i, e := range pe.errSlice {
		ss[i] = e.Error()
	}

	return strings.Join(ss, ",")
}

func NewTypeNotValidError(kind reflect.Kind) TypeNotValidError {
	return TypeNotValidError{kind: kind}
}

func (e TypeNotValidError) Error() string {
	return fmt.Sprintf("Kind must be a struct or pointer to struct, got %q", e.kind)
}
