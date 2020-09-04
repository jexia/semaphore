package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

var errNotAnObject = errors.New("not an object")

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
