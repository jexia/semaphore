package jwt

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Claims interface.
type Claims interface {
	jwt.Claims

	Subject() string
}

// Reader is a claims reader interface.
type Reader interface {
	Read(token string, recv Claims) error
}

// ReaderFunc allows regular function with a proper signature to be used as a Reader.
type ReaderFunc func(string, jwt.Claims) error

func (fn ReaderFunc) Read(token string, recv Claims) error { return fn(token, recv) }

func getAuthorizartionValue(value interface{}) (string, error) {
	authValue, is := value.(string)
	if !is {
		return "", errInvalidValueType
	}

	parts := strings.Split(authValue, " ")
	if len(parts) != 2 {
		return "", errMalformedAuthValue
	}

	if kind := strings.ToLower(strings.TrimSpace(parts[0])); kind != "bearer" {
		return "", errUnsupportedAuthMethod{kind: kind}
	}

	return parts[1], nil
}
