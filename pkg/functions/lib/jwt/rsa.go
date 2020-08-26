package jwt

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
)

// RSA creates a JWT validator with provided public key.
func RSA(publicKey *rsa.PublicKey) ReaderFunc {
	return func(tokenString string, recv jwt.Claims) error {
		token, err := jwt.ParseWithClaims(tokenString, recv, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errUnexpectedSigningMethod{
					actual:   token.Header["alg"],
					expected: "RS265/384/512/...",
				}
			}

			return publicKey, nil
		})

		if err != nil {
			return err
		}

		if !token.Valid {
			return errInvalidToken
		}

		return nil
	}
}
