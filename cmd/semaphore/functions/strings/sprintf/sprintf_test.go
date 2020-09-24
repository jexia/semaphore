package sprintf

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestFunction(t *testing.T) {
	t.Run("should return an error when no arguments provided", func(t *testing.T) {
		var _, _, err = Function()

		if !errors.Is(err, errNoArguments) {
			t.Errorf("error '%s' was expected", errNoArguments)
		}
	})

	t.Run("should return an error when the format has invalid type", func(t *testing.T) {
		var _, _, err = Function(&specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.Bool,
				},
			},
		})

		if !errors.Is(err, errInvalidFormat) {
			t.Errorf("unexpected error '%s', expected '%s'", err, errInvalidFormat)
		}
	})

	t.Run("should return an error when the format is a reference", func(t *testing.T) {
		var _, _, err = Function(&specs.Property{
			Template: specs.Template{
				Reference: &specs.PropertyReference{
					Resource: template.InputResource,
					Path:     "reference",
				},
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		})

		if !errors.Is(err, errNoReferenceSupport) {
			t.Errorf("unexpected error '%s', expected '%s'", err, errNoReferenceSupport)
		}
	})

	t.Run("should return an error when the format is missing", func(t *testing.T) {
		var _, _, err = Function(&specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		})

		if !errors.Is(err, errNoFormat) {
			t.Errorf("unexpected error '%s', expected '%s'", err, errNoFormat)
		}
	})

	t.Run("should propagate scanner error", func(t *testing.T) {
		var _, _, err = Function(&specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Default: "%z",
				},
			},
		})

		if !errors.Is(err, errUnknownFormatter) {
			t.Errorf("unexpected error '%s', expected '%s'", err, errUnknownFormatter)
		}
	})

	t.Run("should return an error when the number of arguments does not match the number of verbs", func(t *testing.T) {
		var _, _, err = Function(&specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Default: "%s",
				},
			},
		})

		if !errors.As(err, &errInvalidArguments{}) {
			t.Errorf("unexpected error '%s', expected '%T'", err, errInvalidArguments{})
		}
	})

	t.Run("should return an error when the type of provided argument does not match expected", func(t *testing.T) {
		var _, _, err = Function(
			&specs.Property{
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type:    types.String,
						Default: "%s",
					},
				},
			},
			&specs.Property{
				Name: "foo",
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.Bool,
					},
				},
			},
		)

		if !errors.As(err, &errCannotFormat{}) {
			t.Errorf("unexpected error '%s', expected '%T'", err, errInvalidArguments{})
		}
	})

	t.Run("should return outputs and executable with no error", func(t *testing.T) {
		var (
			outputs, executable, err = Function(
				&specs.Property{
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.String,
							Default: "develop %s, not hard!",
						},
					},
				},
				&specs.Property{
					Name: "foo",
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			)

			expected = &specs.Property{
				Name:  "sprintf",
				Label: labels.Optional,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.String,
					},
				},
			}
		)

		if err != nil {
			t.Errorf("unexpected error '%s'", err)
		}

		if !reflect.DeepEqual(outputs, expected) {
			t.Errorf("invalid outputs '%v', expected '%v'", outputs, expected)
		}

		if executable == nil {
			t.Error("executable was not expected to be 'nil'")

			return
		}

		var refs = references.NewReferenceStore(1)

		executable(refs)
	})
}
