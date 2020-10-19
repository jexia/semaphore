package jwt

import (
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

const (
	paramClaims = "claims"
	propSubject = "subject"
)

func outputs() *specs.Property {
	return &specs.Property{
		Label: labels.Required,
		Template: specs.Template{
			Message: specs.Message{
				propSubject: {
					Name:  propSubject,
					Path:  propSubject,
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			},
		},
	}
}

// New creates custom function to validate JWT with provided validator.
func New(reader Reader, newClaims func() Claims) functions.Intermediate {
	return func(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
		if len(args) != 1 {
			return nil, nil, errInvalidNumberOfArguments{
				actual:   len(args),
				expected: 1,
			}
		}

		if args[0].Type() != types.String {
			return nil, nil, errInvalidArgumentType{
				actual:   args[0].Type(),
				expected: types.String,
			}
		}

		return outputs(), executable(reader, args[0], newClaims), nil
	}
}

func executable(reader Reader, token *specs.Property, newClaims func() Claims) func(store references.Store) error {
	return func(store references.Store) error {
		var value = token.DefaultValue()

		if token.Reference != nil {
			ref := store.Load(token.Reference.Resource, token.Reference.Path)
			if ref != nil {
				value = ref.Value
			}
		}

		authValue, err := getAuthorizartionValue(value)
		if err != nil {
			return err
		}

		var claimsObj = newClaims()

		if err := reader.Read(authValue, claimsObj); err != nil {
			return err
		}

		store.StoreValue(paramClaims, propSubject, claimsObj.Subject())

		return nil
	}
}
