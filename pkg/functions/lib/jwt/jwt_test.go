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
		token = &specs.Property{
			Name:  "claims",
			Path:  "claims",
			Label: labels.Required,
			Template: specs.Template{
				Reference: &specs.PropertyReference{
					Resource: "input",
					Path:     "authorization",
				},
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		}

		store = references.NewReferenceStore(1)
	)

	t.Run("should propagate Reader error", func(t *testing.T) {
		store.StoreReference("input", &references.Reference{
			Path:  "authorization",
			Value: "Bearer expected.jwt",
		})

		var (
			errExpected = errors.New("expected error")

			reader = ReaderFunc(func(token string, recv jwt.Claims) error { return errExpected })

			fn = executable(reader, token, func() Claims { return new(claims) })
		)

		if err := fn(store); !errors.Is(err, errExpected) {
			t.Errorf("uexpected error: %s", err)
		}
	})

	t.Run("should return an error when unable to get authorization value", func(t *testing.T) {
		store.StoreReference("input", &references.Reference{
			Path:  "authorization",
			Value: "invalid.jwt",
		})

		var (
			reader = ReaderFunc(func(token string, recv jwt.Claims) error { return nil })

			fn = executable(reader, token, func() Claims { return new(claims) })
		)

		if err := fn(store); !errors.Is(err, errMalformedAuthValue) {
			t.Errorf("uexpected error: %s", err)
		}
	})

	t.Run("should save the subject to the reference store", func(t *testing.T) {
		store.StoreReference("input", &references.Reference{
			Path:  "authorization",
			Value: "Bearer expected.jwt",
		})

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

			fn = executable(reader, token, func() Claims { return new(claims) })
		)

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

func TestNew(t *testing.T) {
	var fn = New(nil, nil)

	t.Run("should return an error when invalid number of argumets provided", func(t *testing.T) {
		_, _, err := fn()

		if !errors.As(err, new(errInvalidNumberOfArguments)) {
			t.Errorf("unexpected error: %T", err)
		}
	})

	t.Run("should return an error providing invalid argument", func(t *testing.T) {
		var (
			arg = &specs.Property{
				Label: labels.Required,
				Template: specs.Template{
					Message: specs.Message{},
				},
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
				Label: labels.Required,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.String,
					},
				},
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
