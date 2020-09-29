package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

var errNotAnObject = errors.New("not an object")

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

// ErrUndefinedSpecs occurs when spacs are nil
type ErrUndefinedSpecs struct{}

// Error returns a description of the given error as a string
func (e ErrUndefinedSpecs) Error() string {
	return fmt.Sprint("no object specs defined")
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedSpecs) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedSpecs",
		Message: e.Error(),
	}
}
