package formatter

import (
	"errors"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// String formatter.
type String struct{}

func (String) String() string { return "s" }

// Format provided argument as a string.
func (String) Format(store references.Store, argument *specs.Property) (string, error) {
	var value interface{}

	if argument.Default != nil {
		value = argument.Default
	}

	if argument.Reference != nil {
		if ref := store.Load(argument.Reference.Resource, argument.Reference.Path); ref != nil {
			value = ref.Value
		}
	}

	// TODO: maybe return an empty string
	if value == nil {
		return "", errors.New("not set")
	}

	casted, ok := value.(string)
	if !ok {
		return "", errors.New("not a string")
	}

	return casted, nil
}
