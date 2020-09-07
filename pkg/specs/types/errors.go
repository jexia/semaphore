package types

import "fmt"

// ErrUnknownType is returned when unable to recognize provided data type.
type ErrUnknownType string

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("unknown data type %q", string(e))
}
