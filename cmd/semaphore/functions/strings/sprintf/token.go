package sprintf

import (
	"fmt"
)

type Token interface {
	fmt.Stringer

	isTokenInterface()
}

func countVerbs(tokens []Token) (total int) {
	for _, token := range tokens {
		if _, ok := token.(Verb); ok {
			total++
		}
	}

	return
}

type Constant string

func (c Constant) isTokenInterface() {}

func (c Constant) String() string { return string(c) }

type Verb struct{ formatter Formatter }

func (v Verb) isTokenInterface() {}

func (v Verb) String() string { return "%" + v.formatter.String() }
