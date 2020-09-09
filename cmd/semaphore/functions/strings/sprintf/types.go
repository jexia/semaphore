package sprintf

import (
	"fmt"
)

type Formatter interface {
	fmt.Stringer

	Format(v interface{}) (string, error)
}

// type Validator interface {
// 	Validate(args ...*specs.Property) error
// }
//
// type Builder interface {
// 	Build() string
// }
