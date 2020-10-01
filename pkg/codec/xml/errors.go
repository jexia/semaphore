package xml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

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

type errFailedToEncodeProperty struct {
	property string
	inner    error
}

func (e errFailedToEncodeProperty) Unwrap() error { return e.inner }

func (e errFailedToEncodeProperty) Error() string {
	return fmt.Sprintf("failed to encode property '%s': %s", e.property, e.inner)
}

type errFailedToDecodeProperty struct {
	property string
	inner    error
}

func (e errFailedToDecodeProperty) Unwrap() error { return e.inner }

func (e errFailedToDecodeProperty) Error() string {
	return fmt.Sprintf("failed to decode property '%s': %s", e.property, e.inner)
}
