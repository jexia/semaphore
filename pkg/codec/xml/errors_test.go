package xml

import (
	"encoding/xml"
	"testing"
)

func TestErrUndefinedProperty(t *testing.T) {
	var (
		err      = errUndefinedProperty("foo")
		expected = `undefined property "foo"`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}

func TestErrUnknownEnum(t *testing.T) {
	var (
		err      = errUnknownEnum("pending")
		expected = `unrecognized enum value "pending"`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}

func TestErrUnexpectedToken(t *testing.T) {
	var (
		err = errUnexpectedToken{
			actual: xml.CharData{},
			expected: []xml.Token{
				xml.StartElement{},
				xml.EndElement{},
			},
		}
		expected = `unexpected element "xml.CharData", expected one of ["xml.StartElement", "xml.EndElement"]`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}
