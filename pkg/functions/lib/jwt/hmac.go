package jwt

import "github.com/dgrijalva/jwt-go"

// HMAC creates simple JWT validator with provided secret string.
func HMAC(secretString string) ReaderFunc {
	return func(tokenString string, recv jwt.Claims) error {
		token, err := jwt.ParseWithClaims(tokenString, recv, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errUnexpectedSigningMethod{
					actual:   token.Header["alg"],
					expected: "HS265/384/512",
				}
			}

			return []byte(secretString), nil
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
