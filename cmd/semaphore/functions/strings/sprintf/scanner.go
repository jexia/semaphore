package sprintf

import (
	"strconv"
)

// Scanner scans provided string splitting it up by tokens.
type Scanner interface {
	Scan(input string) ([]Token, error)
}

type state func(prev Token, input string, start int) (token Token, pos int, next state, err error)

// defaultScanner is a default implementation of Scanner interface.
type defaultScanner struct {
	FormatterDetector
}

// NewScanner creates a stateful scanner with provided formatter detector.
// Note that the instance can be used only once.
func NewScanner(detector FormatterDetector) Scanner {
	return &defaultScanner{
		FormatterDetector: detector,
	}
}

// Scan the input for tokens.
func (s *defaultScanner) Scan(input string) ([]Token, error) {
	var (
		pos    = 0
		next   = s.scanConstant
		tokens []Token
		token  Token
	)

	for pos < len(input) {
		var err error

		token, pos, next, err = next(token, input, pos)

		if err != nil {
			// TODO: use the position to provide more info
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (s *defaultScanner) scanConstant(_ Token, input string, start int) (Token, int, state, error) {
	var curr = start

	for ; curr < len(input); curr++ {
		if rune(input[curr]) == '%' {
			break
		}
	}

	return Constant(input[start:curr]), curr + 1, s.scanPrecision, nil
}

func (s *defaultScanner) scanVerb(prev Token, input string, start int) (Token, int, state, error) {
	constructor, ok := s.Detect(input[start:])
	if !ok {
		return nil, start, nil, errUnknownFormatter
	}

	formatter, err := constructor.Formatter(prev.(Precision))
	if err != nil {
		return nil, start, nil, err
	}

	var token = Verb{
		Verb:      constructor.String(),
		Formatter: formatter,
	}

	return token, start + len(constructor.String()), s.scanConstant, nil
}

func (s *defaultScanner) scanPrecision(_ Token, input string, start int) (Token, int, state, error) {
	const (
		waitWidth = iota
		waitScale
		allDone
	)

	var (
		state     = waitWidth
		curr      = start
		precision Precision

		width = []byte{'0'}
		scale = []byte{'0'}
	)

LOOP:
	for ; curr < len(input); curr++ {
		char := rune(input[curr])
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch state {
			case waitWidth:
				width = append(width, byte(char))
			case waitScale:
				scale = append(scale, byte(char))
			}
		case '.':
			switch state {
			case waitWidth:
				state = waitScale

				continue
			case waitScale:
				return nil, curr, nil, errMalformedPrecision
			}
			state = waitScale
		default:
			int64Width, err := strconv.ParseInt(string(width), 10, 64)
			if err != nil {
				return nil, curr, nil, err
			}

			precision.Width = int(int64Width)

			int64Scale, err := strconv.ParseInt(string(scale), 10, 64)
			if err != nil {
				return nil, curr, nil, err
			}

			precision.Scale = int(int64Scale)

			state = allDone

			break LOOP
		}
	}

	if state != allDone {
		return nil, curr, nil, errMissingFormat
	}

	return precision, curr, s.scanVerb, nil
}
