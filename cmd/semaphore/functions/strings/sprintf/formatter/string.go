package formatter

import "fmt"

type String struct{}

func (String) String() string { return "s" }

func (String) Format(store references.Store, v interface{}) (string, error) {
	var value interface{}

	if argument.Default != nil {
		value = argument.Default
	}

	if argument.Reference != nil {
		if ref := store.Load(argument.Reference.Resource, argument.Reference.Path); ref != nil {
			value = ref.Value
		}
	}

	return fmt.Sprint(v), nil
}
