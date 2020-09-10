package sprintf

import (
	"fmt"
	"strconv"
)

// Token represents one of the available tokens.
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

func onlyVerbs(tokens []Token) (verbs []Verb) {
	for _, token := range tokens {
		if verb, ok := token.(Verb); ok {
			verbs = append(verbs, verb)
		}
	}

	return
}

// Constant is a token that does not need any formatting.
type Constant string

func (Constant) isTokenInterface() {}

func (c Constant) String() string { return string(c) }

// Verb is a token that is used to print a single input argument.
type Verb struct {
	Verb      string
	Formatter Formatter
}

func (Verb) isTokenInterface() {}

func (v Verb) String() string { return v.Verb }

// Precision is a token which describes the precision and used to create a verb.
type Precision struct {
	Width int64
	Scale int64
}

func (Precision) isTokenInterface() {}

func (p Precision) String() string {
	return "%" + strconv.FormatInt(p.Width, 10) + "." + strconv.FormatInt(p.Scale, 10)
}
