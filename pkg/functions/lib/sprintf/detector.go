package sprintf

import iradix "github.com/hashicorp/go-immutable-radix"

// FormatterDetector detects formatter (if possible) from provided input.
type FormatterDetector interface {
	Detect(input string) (Constructor, bool)
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
func (r *Radix) Register(constructor Constructor) error {
	var ok bool

	r.tree, _, ok = r.tree.Insert([]byte(constructor.String()), constructor)
	if ok {
		return errVerbConflict{constructor}
	}

	return nil
}

// Detect if input string starts with one of the registered verbs.
func (r *Radix) Detect(input string) (Constructor, bool) {
	_, v, ok := r.tree.Root().LongestPrefix([]byte(input))
	if !ok {
		return nil, false
	}

	return v.(Constructor), true
}
