package jwt

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

var (
	errInvalidToken       = errors.New("token is invalid")
	errInvalidValueType   = errors.New("invalid value, expected a string")
	errMalformedAuthValue = errors.New("malformed auhtorization value")
)

type errUnsupportedAuthMethod struct {
	kind string
}

func (e errUnsupportedAuthMethod) Error() string {
	return fmt.Sprintf("unsuported authorization method (%s)", e.kind)
}

type errInvalidNumberOfArguments struct {
	actual, expected int
}

func (e errInvalidNumberOfArguments) Error() string {
	return fmt.Sprintf("invalid number of arguments (%d), expected (%d)", e.actual, e.expected)
}

type errInvalidArgumentType struct {
	actual, expected types.Type
}

func (e errInvalidArgumentType) Error() string {
	return fmt.Sprintf("invalid argument type (%s), expected (%s)", e.actual, e.expected)
}

type errUnexpectedSigningMethod struct {
	actual, expected interface{}
}

func (e errUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method (%v), expected (%v)", e.actual, e.expected)
}
