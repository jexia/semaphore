package formencoded

import (
	"errors"
	"fmt"
)

var errNilSchema = errors.New("nil schema")

type errUndefinedProperty string

func (e errUndefinedProperty) Error() string {
	return fmt.Sprintf("undefined property %q", string(e))
}

type errUnknownLabel string

func (e errUnknownLabel) Error() string {
	return fmt.Sprintf("unknown label %q", string(e))
}

type errUnknownType string

func (e errUnknownType) Error() string {
	return fmt.Sprintf("unknown data type %q", string(e))
}
