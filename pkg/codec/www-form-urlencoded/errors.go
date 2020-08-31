package formencoded

import "fmt"

type errUndefinedProperty string

func (e errUndefinedProperty) Error() string {
	return fmt.Sprintf("undefined property %q", string(e))
}

type errUnknownLabel string

func (e errUnknownLabel) Error() string {
	return fmt.Sprintf("unknown label %q", string(e))
}
