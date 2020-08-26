package jwt

import (
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type claims struct{ jwt.StandardClaims }

func (c claims) Subject() string { return c.StandardClaims.Subject }

func TestExecutable(t *testing.T) {
	var (
		reader = ReaderFunc(func(token string, recv jwt.Claims) error {
			if token != "expected.jwt" {
				return errors.New("unexpected token")
			}

			casted, ok := recv.(*claims)
			if !ok {
				return errors.New("invalid receiver")
			}

			casted.StandardClaims.Subject = "expected subject"

			return nil
		})

		token = &specs.Property{
			Name: "claims",
			Path: "claims",
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "authorization",
			},
			Type:  types.String,
			Label: labels.Required,
		}

		store = references.NewReferenceStore(1)

		fn = executable(reader, token, func() Claims { return new(claims) })
	)

	store.StoreReference("input", &references.Reference{
		Path:  "authorization",
		Value: "Bearer expected.jwt",
	})

	if err := fn(store); err != nil {
		t.Errorf("uexpected error: %s", err)
	}

	subject := store.Load(paramClaims, propSubject)
	if subject == nil {
		t.Error("subject must be stored")
	}

	if subject.Value != "expected subject" {
		t.Errorf("unexpected subject: %v", subject.Value)
	}
}
