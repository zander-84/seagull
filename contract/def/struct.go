package def

// E represents a  element .
type E struct {
	Key   string
	Value interface{}
}

// D is an ordered representation of a  document
type D []E

type M map[string]interface{}

// A represents a array.
type A []any

// K represents a key.
type K struct {
	Key   string
	Alias []any
}
