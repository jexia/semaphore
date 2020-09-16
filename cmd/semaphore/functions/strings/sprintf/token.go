package sprintf

import (
	"fmt"
	"strconv"
)

// Token represents one of the available tokens.
type Token interface {
	fmt.Stringer
}

// Constant is a token that does not need any formatting.
type Constant string

func (c Constant) String() string { return string(c) }

// Verb is a token that is used to print a single input argument.
type Verb struct {
	TypeChecker
	Verb      string
	Formatter Formatter
}

func (v Verb) String() string { return v.Verb }

// Precision is a token which describes the precision and used to create a verb.
type Precision struct {
	Width int64
	Scale int64
}

func (p Precision) String() string {
	return "%" + strconv.FormatInt(p.Width, 10) + "." + strconv.FormatInt(p.Scale, 10)
}
