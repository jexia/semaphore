package specs

import (
	"errors"
	"fmt"
)

var (
	errNilTemplate = errors.New("provided template is nil")
)

type errTypeMismatch struct {
	Inner error
}

func (e errTypeMismatch) Error() string {
	return fmt.Sprintf("type mismatch: %s", e.Inner)
}

func (e errTypeMismatch) Unwrap() error {
	return e.Inner
}
