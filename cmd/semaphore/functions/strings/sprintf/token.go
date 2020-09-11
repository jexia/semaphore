package sprintf

import (
	"fmt"
	"strconv"
)

const (
	// TConstant is a constant token which should be printed as is.
	TConstant Kind = iota
	// TPrecision is a token that is not printed itself but is a part of verb.
	TPrecision
	// TVerb is a token which contains the format verb.
	TVerb
)

// Kind represents token kind.
type Kind int

// Token represents one of the available tokens.
type Token interface {
	fmt.Stringer

	Kind() Kind
}

// Constant is a token that does not need any formatting.
type Constant string

// Kind returns token kind.
func (Constant) Kind() Kind { return TConstant }

func (c Constant) String() string { return string(c) }

// Verb is a token that is used to print a single input argument.
type Verb struct {
	TypeChecker
	Verb      string
	Formatter Formatter
}

// Kind returns token kind.
func (Verb) Kind() Kind { return TVerb }

func (v Verb) String() string { return v.Verb }

// Precision is a token which describes the precision and used to create a verb.
type Precision struct {
	Width int64
	Scale int64
}

// Kind returns token kind.
func (Precision) Kind() Kind { return TPrecision }

func (p Precision) String() string {
	return "%" + strconv.FormatInt(p.Width, 10) + "." + strconv.FormatInt(p.Scale, 10)
}
