package jwt

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestErrUnsupportedAuthMethod(t *testing.T) {
	var (
		expected = "unsuported authorization method (unknown)"
		err      = errUnsupportedAuthMethod{kind: "unknown"}
	)

	if msg := err.Error(); msg != expected {
		t.Errorf("invalid error message %q, expected %q", msg, expected)
	}
}

func TestErrInvalidNumberOfArguments(t *testing.T) {
	var (
		expected = "invalid number of arguments (1), expected (2)"
		err      = errInvalidNumberOfArguments{actual: 1, expected: 2}
	)

	if msg := err.Error(); msg != expected {
		t.Errorf("invalid error message %q, expected %q", msg, expected)
	}
}

func TestErrInvalidArgumentType(t *testing.T) {
	var (
		expected = "invalid argument type (message), expected (string)"
		err      = errInvalidArgumentType{actual: types.Message, expected: types.String}
	)

	if msg := err.Error(); msg != expected {
		t.Errorf("invalid error message %q, expected %q", msg, expected)
	}
}

func TestErrUnexpectedSigningMethod(t *testing.T) {
	var (
		expected = "unexpected signing method (HMAC256), expected (RSA256)"
		err      = errUnexpectedSigningMethod{actual: "HMAC256", expected: "RSA256"}
	)

	if msg := err.Error(); msg != expected {
		t.Errorf("invalid error message %q, expected %q", msg, expected)
	}
}
