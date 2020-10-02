package xml

import (
	"encoding/xml"
	"errors"
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

func TestErrFailedToEncode(t *testing.T) {
	var (
		err = errFailedToEncode{
			errStack: errStack{
				property: "root",
				inner: errFailedToEncode{
					errStack: errStack{
						property: "nested",
						inner: errFailedToEncode{
							errStack: errStack{
								property: "integer",
								inner:    errors.New("internal error"),
							},
						},
					},
				},
			},
		}
		expected = `failed to encode element: path 'root.nested.integer': internal error`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}

func TestErrFailedToDecode(t *testing.T) {
	var (
		err = errFailedToDecode{
			errStack: errStack{
				property: "root",
				inner: errFailedToDecode{
					errStack: errStack{
						property: "nested",
						inner: errFailedToDecode{
							errStack: errStack{
								property: "integer",
								inner:    errors.New("internal error"),
							},
						},
					},
				},
			},
		}
		expected = `failed to decode element: path 'root.nested.integer': internal error`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}
