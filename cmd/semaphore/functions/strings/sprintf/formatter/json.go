package formatter

import (
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// JSON formatter.
type JSON struct{}

func (JSON) String() string { return "json" }

// Format prints provided argument in a JSON format.
func (JSON) Format(store references.Store, argument *specs.Property) (string, error) {
	// data, err := json.Marshal(v)
	// if err != nil {
	// 	return "", err
	// }
	//
	// return string(data), nil

	panic("not implemented")
}
