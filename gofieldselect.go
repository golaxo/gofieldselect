package gofieldselect

import (
	"reflect"
	"strings"

	"github.com/golaxo/gofieldselect/internal/lexer"
)

// GetWithReflection creates a new instance [T] with only the fields specified in [n].
// It goes, using reflection, through all the fields in the type [T] and if the field is exported
// and by either checking the JSON tag or the field name, setting a default value or the
// value that comes from the source.
//
//nolint:exhaustive,gocognit,nestif,funlen // refactor later
func GetWithReflection[T any](n Node, source T) (T, error) {
	var zero T

	rv := reflect.ValueOf(source)
	rt := rv.Type()

	// Helper to apply selection from src to dst (both struct values)
	var apply func(Node, reflect.Value, reflect.Value) error

	apply = func(node Node, src, dst reflect.Value) error {
		// If all identifiers, copy the whole struct
		if _, ok := node.(AllIdentifiers); ok {
			dst.Set(src)

			return nil
		}

		st := src.Type()
		for i := range st.NumField() {
			sf := st.Field(i)

			dv := dst.Field(i)
			if !dv.CanSet() {
				continue
			}

			// Determine a selection key: JSON tag (name part) or field name
			name := sf.Name
			if tag := sf.Tag.Get("json"); tag != "" {
				tagName := tag
				if idx := strings.Index(tagName, ","); idx >= 0 {
					tagName = tagName[:idx]
				}

				if tagName == "-" {
					// Unexported for JSON selection; fall back to skipping
					continue
				}

				if tagName != "" { // explicit empty means use field name
					name = tagName
				}
			}

			ident, ok := node.SelectField(name)
			if !ok {
				continue
			}

			sv := src.Field(i)

			// If there is a child selection, apply recursively for structs
			if ident.Child != nil {
				switch sv.Kind() {
				case reflect.Struct:
					// Recurse into struct
					if err := apply(ident.Child, sv, dv); err != nil {
						return err
					}
				case reflect.Ptr:
					if sv.IsNil() {
						// source is nil; leave destination as zero (nil)
						continue
					}

					if sv.Elem().Kind() == reflect.Struct {
						if dv.IsNil() {
							dv.Set(reflect.New(sv.Elem().Type()))
						}

						if err := apply(ident.Child, sv.Elem(), dv.Elem()); err != nil {
							return err
						}
					} else {
						// Non-struct pointer: copy as is
						dv.Set(sv)
					}
				default:
					// Non-struct: copy value
					dv.Set(sv)
				}

				continue
			}

			// No child selection: copy field value directly
			dv.Set(sv)
		}

		return nil
	}

	switch rt.Kind() {
	case reflect.Ptr:
		// Expect pointer to struct
		if rt.Elem().Kind() != reflect.Struct {
			return zero, NewTypeNotValidError(rt.Kind())
		}
		// Create a new pointer and apply
		dst := reflect.New(rt.Elem())
		// If the node is AllIdentifiers, copy the whole value
		if _, ok := n.(AllIdentifiers); ok {
			if !rv.IsNil() {
				dst.Elem().Set(rv.Elem())
			}

			//nolint:errcheck // it's always T
			return dst.Interface().(T), nil
		}

		if rv.IsNil() {
			// Source is nil; return a new zero pointer (nil) because copying selected fields from nil is undefined.
			// Keep it nil to be consistent with zero value behavior.
			var nilPtr T

			return nilPtr, nil
		}

		if err := apply(n, rv.Elem(), dst.Elem()); err != nil {
			return zero, err
		}

		//nolint:errcheck // it's always T
		return dst.Interface().(T), nil
	case reflect.Struct:
		dst := reflect.New(rt).Elem()
		if err := apply(n, rv, dst); err != nil {
			return zero, err
		}

		//nolint:errcheck // it's always T
		return dst.Interface().(T), nil
	default:
		return zero, NewTypeNotValidError(rt.Kind())
	}
}

func Get[T any](n Node, fieldName string, originalValue T) T {
	_, ok := n.SelectField(fieldName)
	if !ok {
		// return default value, for ptr nil, for non pointers zero value
		var zero T

		return zero
	}

	return originalValue
}

func Parse(fieldSelection string) (Node, error) {
	p := newParser(lexer.New(fieldSelection))

	n := p.parse()
	if len(p.Errors()) > 0 {
		return nil, NewParsingError(p.Errors())
	}

	return n, nil
}
