package sprintf

import (
	"fmt"
	"strings"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

type Printer interface {
	Print(store references.Store, args ...*specs.Property) (string, error)
}

type defaultPrinter struct {
	tokens []Token
}

func NewPrinter(tokens []Token) Printer {
	return &defaultPrinter{
		tokens: tokens,
	}
}

func (p *defaultPrinter) Print(store references.Store, args ...*specs.Property) (string, error) {
	var (
		verbPos int
		builder strings.Builder
	)

	for _, token := range p.tokens {
		switch t := token.(type) {
		case Constant:
			if _, err := builder.WriteString(string(t)); err != nil {
				return "", err
			}
		case Verb:
			str, err := t.Print(store, args[verbPos])
			if err != nil {
				return "", err
			}

			if _, err := builder.WriteString(str); err != nil {
				return "", err
			}

			verbPos++
		default:
			return "", fmt.Errorf("unexpected token %T", t)
		}
	}

	return builder.String(), nil
}
