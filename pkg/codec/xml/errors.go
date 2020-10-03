package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

var errNoSchema = errors.New("no object specs defined")

type errUndefinedProperty string

func (e errUndefinedProperty) Error() string {
	return fmt.Sprintf("undefined property %q", string(e))
}

type errUnknownEnum string

func (e errUnknownEnum) Error() string {
	return fmt.Sprintf("unrecognized enum value %q", string(e))
}

type errUnexpectedToken struct {
	actual   xml.Token
	expected []xml.Token
}

func (e errUnexpectedToken) printExpected() string {
	var builder strings.Builder

	for index, token := range e.expected {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf(`"%T"`, token))
	}

	return builder.String()
}

func (e errUnexpectedToken) Error() string {
	return fmt.Sprintf(`unexpected element "%T", expected one of [%s]`, e.actual, e.printExpected())
}

type nestedError interface {
	path() string
	unwrap() error
}

type errStack struct {
	property string
	inner    error
}

func (e errStack) Unwrap() error { return e.inner }

func (e errStack) path() string {
	casted, ok := e.inner.(nestedError)
	if !ok {
		return e.property
	}

	return buildPath(e.property, casted.path())
}

func (e errStack) unwrap() error {
	caseted, ok := e.inner.(nestedError)
	if !ok {
		return e.inner
	}

	return caseted.unwrap()
}

type errFailedToEncode struct{ errStack }

func (e errFailedToEncode) Error() string {
	return fmt.Sprintf("failed to encode element: path '%s': %s", e.path(), e.unwrap())
}

type errFailedToDecode struct{ errStack }

func (e errFailedToDecode) Unwrap() error { return e.inner }

func (e errFailedToDecode) Error() string {
	return fmt.Sprintf("failed to decode element: path '%s': %s", e.path(), e.unwrap())
}
