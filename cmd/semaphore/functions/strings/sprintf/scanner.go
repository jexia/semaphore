package sprintf

import (
	"fmt"
)

// Scanner scans provided string splitting it up by tokens.
type Scanner interface {
	Scan(input string) ([]Token, error)
}

type state func(input string, start int) (token Token, pos int, next state, err error)

// DefaultScanner is a default implementation of Scanner interface.
type DefaultScanner struct {
	FormatterDetector
}

// NewDefaultScanner creates a default scanner with provided formatter detector.
func NewDefaultScanner(detector FormatterDetector) *DefaultScanner {
	return &DefaultScanner{
		FormatterDetector: detector,
	}
}

// Scan the input for tokens.
func (s *DefaultScanner) Scan(input string) ([]Token, error) {
	var (
		pos    = 0
		next   = s.scanConstant
		tokens []Token
	)

	for pos < len(input) {
		var (
			token Token
			err   error
		)

		token, pos, next, err = next(input, pos)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (s *DefaultScanner) scanVerb(input string, start int) (Token, int, state, error) {
	formatter, ok := s.Detect(input[start:])
	if !ok {
		return nil, 0, nil, fmt.Errorf("unknown formatter")
	}

	var verb = Verb{
		formatter: formatter,
	}

	return verb, start + len(formatter.String()), s.scanConstant, nil
}

func (s *DefaultScanner) scanConstant(input string, start int) (Token, int, state, error) {
	var curr = start

	for ; curr < len(input); curr++ {
		if rune(input[curr]) == '%' {
			break
		}
	}

	return Constant(input[start:curr]), curr + 1, s.scanVerb, nil
}
