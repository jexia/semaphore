package sprintf

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

type Token interface {
	fmt.Stringer

	isTokenInterface()
}

func countVerbs(tokens []Token) (total int) {
	for _, token := range tokens {
		if _, ok := token.(Verb); ok {
			total++
		}
	}

	return
}

type Constant string

func (c Constant) isTokenInterface() {}

func (c Constant) String() string { return string(c) }

type Verb struct{ formatter Formatter }

func (v Verb) isTokenInterface() {}

func (v Verb) String() string { return "%" + v.formatter.String() }

func (v Verb) Print(store references.Store, argument *specs.Property) (string, error) {
	var value interface{}

	if argument.Default != nil {
		value = argument.Default
	}

	if argument.Reference != nil {
		if ref := store.Load(argument.Reference.Resource, argument.Reference.Path); ref != nil {
			value = ref.Value
		}
	}

	return v.formatter.Format(value)
}
