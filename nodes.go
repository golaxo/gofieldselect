package gofieldselect

var (
	_ Node = new(Identifiers)
	_ Node = new(AllIdentifiers)
)

type (
	// Node is the top level of the field selection.
	Node interface {
		// SelectField method to indicate if the field is selected.
		SelectField(fieldName string) (Identifier, bool)
		node()
	}

	// Identifiers is the node with the list of fields that should be selected.
	Identifiers []Identifier

	// AllIdentifiers is the Node to indicate that every field should be selected.
	AllIdentifiers struct{}

	// Identifier holder that indicates the JSON field name and the children selection for that field
	// E.g. `name` or `address(street,number)`.
	Identifier struct {
		Value string
		Child Node
	}
)

func (is Identifiers) SelectField(fieldName string) (Identifier, bool) {
	for _, i := range is {
		if i.Value == fieldName {
			return i, true
		}
	}

	return Identifier{}, false
}

func (is Identifiers) node() {}

func (a AllIdentifiers) SelectField(fieldName string) (Identifier, bool) {
	return Identifier{
		Value: fieldName,
		Child: AllIdentifiers{},
	}, true
}

func (a AllIdentifiers) node() {}
