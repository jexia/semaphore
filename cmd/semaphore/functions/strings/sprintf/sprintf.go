package sprintf

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func sprintfOutputs() *specs.Property {
	return &specs.Property{
		Name:  "sprintf",
		Type:  types.String,
		Label: labels.Optional,
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

// Sprintf formats and returns a string without printing it anywhere.
func Sprintf(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
	if actual := len(args); actual < 1 {
		return nil, nil, fmt.Errorf("at least 1 argument is expected, provided %d", actual)
	}

	var format = args[0]

	if format.Type != types.String {
		return nil, nil, errors.New("format must be a string")
	}

	if format.Default == nil {
		return nil, nil, fmt.Errorf("invalid format")
	}

	tokens, err := scanner.Scan(format.Default.(string))
	if err != nil {
		return nil, nil, err
	}

	args = args[1:]

	if actual, expected := len(args), countVerbs(tokens); actual != expected {
		return nil, nil, fmt.Errorf("invalid number of input arguments %d, expected %d", actual, expected)
	}

	// forsindex

	return sprintfOutputs(), sprintfExecutable(NewPrinter(tokens), args...), nil
}
