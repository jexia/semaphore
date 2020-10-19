package jwt

import (
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestHMAC(t *testing.T) {
	var reader = HMAC("secret")

	t.Run("should return an error when the token with uexpected signing method is provided", func(t *testing.T) {
		var (
			token  = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30.Ps-lkT0lKrfdrkHrN-efCNfmGAV-nLmO24xl7VBJMZZxRJes6P2fH9B4wQHRMgmjwvijpmYN4qBuFFtibtEYN0h3_KfO9IXi3FspoFfDCl1C3oVRE_OAsqW6k148TGTTZ28ozlzs2ngwLFJpt9TYmkkr_MOsIPpX6jT00iU5758CPo3Lj714JD8FFKch42Eokcdfbt8b_Rv7TYUaIFJsnrJu77Cei5Em1EFPjD91o58eQsHxMBRiQLODadOf5xrFC4LkOqkdat_l1dSYCGPUEg3PufHTqwDilIZVBujUUZzGE-EXqxumsrtVnU0oHz23tLgpZOfQZqqx2UvOKP63Yg`
			claims = new(claims)
			err    = reader.Read(token, claims)
		)

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
			token  = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MTYyMzkwMjJ9.G3a2cpmqhbZq6Q39-zHH9oRn18fJG_HDIKspoljeKUA`
			claims = new(claims)
			err    = reader.Read(token, claims)
		)

		_, ok := err.(*jwt.ValidationError)
		if !ok {
			t.Errorf("unexpected error type: (%T)", err)
		}
	})

	t.Run("should not fail when the token is valid", func(t *testing.T) {
		var (
			token  = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MTYyMzkwMjJ9.RHTJ6ziBOktHKqiGE-HhBQUrr-7gTJJDdAdg1-r38oI`
			claims = new(claims)
			err    = reader.Read(token, claims)
		)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
