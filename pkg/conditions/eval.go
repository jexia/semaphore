package conditions

import (
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

// Eval evaluates the given condition
func Eval(store refs.Store, condition *specs.Condition) bool {
	ref := store.Load(condition.Reference.Resource, condition.Reference.Path)
	if ref == nil {
		return false
	}

	return true
}
