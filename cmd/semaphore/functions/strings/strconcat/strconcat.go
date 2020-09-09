package strconcat

import (
	"strings"

	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Strconcat compiles the given arguments and constructs a new executable
// function for the given arguments.
func Strconcat(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
	result := &specs.Property{
		Name:  "concat",
		Type:  types.String,
		Label: labels.Optional,
	}

	for _, arg := range args {
		if arg.Type != types.String {
			return nil, nil, ErrInvalidArgument{
				Property: arg,
				Expected: types.String,
				Function: "strconcat",
			}
		}
	}

	handle := func(store references.Store) error {
		result := strings.Builder{}

		for _, arg := range args {
			var value string

			if arg.Default != nil {
				value = arg.Default.(string)
			}

			if arg.Reference != nil {
				ref := store.Load(arg.Reference.Resource, arg.Reference.Path)
				if ref != nil {
					value = ref.Value.(string)
				}
			}

			_, err := result.WriteString(value)
			if err != nil {
				return err
			}
		}

		store.StoreValue("", ".", result.String())
		return nil
	}

	return result, handle, nil
}
