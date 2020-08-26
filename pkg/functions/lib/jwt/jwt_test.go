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
	t.Run("run executable", func(t *testing.T) {
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
	})
}

func TestJWT(t *testing.T) {
	var fn = JWT(nil, nil)

	t.Run("should return an error when invalid number of argumets provided", func(t *testing.T) {
		_, _, err := fn()

		if !errors.As(err, new(errInvalidNumberOfArguments)) {
			t.Errorf("unexpected error: %T", err)
		}
	})

	t.Run("should return an error providing invalid argument", func(t *testing.T) {
		var (
			arg = &specs.Property{
				Type:  types.Message,
				Label: labels.Required,
			}

			_, _, err = fn(arg)
		)

		if !errors.As(err, new(errInvalidArgumentType)) {
			t.Errorf("unexpected error: %T", err)
		}
	})

	t.Run("should return outputs and executable when input argument is valid", func(t *testing.T) {
		var (
			arg = &specs.Property{
				Type:  types.String,
				Label: labels.Required,
			}

			outputs, executable, err = fn(arg)
		)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if outputs == nil {
			t.Error("outputs was expected to be set")
		}

		if executable == nil {
			t.Error("executable was expected to be set")
		}
	})
}
