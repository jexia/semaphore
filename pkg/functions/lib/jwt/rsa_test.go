package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestRSA(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	var reader = RSA(&key.PublicKey)

	t.Run("should return an error when the token with uexpected signing method is provided", func(t *testing.T) {
		var (
			token       = jwt.NewWithClaims(jwt.SigningMethodHS256, new(claims))
			signed, err = token.SignedString([]byte("secret"))
		)

		if err != nil {
			t.Fatal(err)
		}

		err = reader.Read(signed, new(claims))

		errValidation, ok := err.(*jwt.ValidationError)
		if !ok {
			t.Errorf("unexpected error type: (%T)", err)
		}

		if !errors.As(errValidation.Inner, new(errUnexpectedSigningMethod)) {
			t.Errorf("unexpected error: %s", errValidation)
		}
	})

	t.Run("should return an error providing expired token", func(t *testing.T) {
		var (
			token = jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() - 100,
			})
			signed, err = token.SignedString(key)
		)

		if err != nil {
			t.Fatal(err)
		}

		err = reader.Read(signed, new(claims))

		_, ok := err.(*jwt.ValidationError)
		if !ok {
			t.Errorf("unexpected error type: (%T)", err)
		}
	})

	t.Run("should not fail when the token is valid", func(t *testing.T) {
		var (
			token = jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + 100,
			})
			signed, err = token.SignedString(key)
		)

		if err != nil {
			t.Fatal(err)
		}

		err = reader.Read(signed, new(claims))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
