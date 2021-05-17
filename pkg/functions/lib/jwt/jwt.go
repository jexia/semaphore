package jwt

import (
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/labels"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

const (
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
		value := token.DefaultValue()

		if token.Reference != nil {
			ref := store.Load(template.ResourcePath(token.Reference.Resource, token.Reference.Path))
			if ref != nil {
				value = ref.Value
			}
		}

		authValue, err := getAuthorizartionValue(value)
		if err != nil {
			return err
		}

		claimsObj := newClaims()

		if err := reader.Read(authValue, claimsObj); err != nil {
			return err
		}

		store.Store(propSubject, &references.Reference{Value: claimsObj.Subject()})

		return nil
	}
}
