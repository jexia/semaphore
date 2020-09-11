package sprintf

import (
	"strings"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Printer compiles the output using argument values.
type Printer interface {
	Print(store references.Store, args ...*specs.Property) (string, error)
}

// Tokens is a list of tokens (implements Printer interface).
type Tokens []Token

// Print the tokens formatting provided arguments according to the format string.
func (tokens Tokens) Print(store references.Store, args ...*specs.Property) (string, error) {
	var (
		verbPos int
		builder strings.Builder
	)

	for _, token := range tokens {
		switch t := token.(type) {
		case Constant:
			if _, err := builder.WriteString(string(t)); err != nil {
				return "", err
			}
		case Verb:
			str, err := t.Formatter(store, args[verbPos])
			if err != nil {
				return "", err
			}

			if _, err := builder.WriteString(str); err != nil {
				return "", err
			}

			verbPos++
		default:
			// ignore the rest of the tokens
		}
	}

	return builder.String(), nil
}

// Verbs filters and returs tokens of type Verb only.
func (tokens Tokens) Verbs() (verbs []Verb) {
	for _, token := range tokens {
		if verb, ok := token.(Verb); ok {
			verbs = append(verbs, verb)
		}
	}

	return
}
