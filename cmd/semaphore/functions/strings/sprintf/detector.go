package sprintf

import (
	"fmt"

	iradix "github.com/hashicorp/go-immutable-radix"
)

// FormatterDetector detects formatter (if possible) from provided input.
type FormatterDetector interface {
	Detect(input string) (Formatter, bool)
}

// Radix tree based implementation of Formatter registry.
type Radix struct {
	tree *iradix.Tree
}

// NewRadix creates new radix tree based Formatter registry.
func NewRadix() *Radix {
	return &Radix{
		tree: iradix.New(),
	}
}

// Register provided Formatter.
func (r *Radix) Register(formatter Formatter) error {
	var ok bool

	r.tree, _, ok = r.tree.Insert([]byte(formatter.String()), formatter)
	if ok {
		return fmt.Errorf("formatter %q is already registered", formatter)
	}

	return nil
}

// Detect if input string starts with one of the registered Formatters.
func (r *Radix) Detect(input string) (Formatter, bool) {
	_, v, ok := r.tree.Root().LongestPrefix([]byte(input))
	if !ok {
		return nil, false
	}

	return v.(Formatter), true
}
