package sprintf

import (
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func sprintfOutputs() *specs.Property {
	return &specs.Property{
		Name:  "sprintf",
		Label: labels.Optional,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.String,
			},
		},
	}
}

func sprintfExecutable(printer Printer, args ...*specs.Property) func(store references.Store) error {
	return func(store references.Store) error {
		result, err := printer.Print(store, args...)
		if err != nil {
			return err
		}

		store.StoreValue("", ".", result)

		return nil
	}
}

// Function formats and returns a string without printing it anywhere.
func Function(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
	if len(args) < 1 {
		return nil, nil, errNoArguments
	}

	var format = args[0]

	if format.Type() != types.String {
		return nil, nil, errInvalidFormat
	}

	if format.Reference != nil {
		return nil, nil, errNoReferenceSupport
	}

	if format.Scalar.Default == nil {
		return nil, nil, errNoFormat
	}

	tokens, err := scanner.Scan(format.Scalar.Default.(string))
	if err != nil {
		return nil, nil, err
	}

	args = args[1:]

	var (
		printer = Tokens(tokens)
		verbs   = printer.Verbs()
	)

	if actual, expected := len(args), len(verbs); actual != expected {
		return nil, nil, errInvalidArguments{
			actual:   actual,
			expected: expected,
		}
	}

	for index, verb := range printer.Verbs() {
		if !verb.CanFormat(args[index].Type()) {
			return nil, nil, errCannotFormat{
				formatter: verb,
				argument:  args[index],
			}
		}
	}

	return sprintfOutputs(), sprintfExecutable(printer, args...), nil
}
