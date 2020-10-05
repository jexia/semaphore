package prettyerr

import (
	"errors"
)

const (
	// GenericErrorCode is for errors which does not implement Prettifier and cannot be recognized with a default strategy
	GenericErrorCode = "GenericError"
)

// Error describes the error with details, organized in a standard structure.
// Error also implements `error` interface.
type Error struct {
	// Original error, might be empty.
	Original error `json:"-"`
	// Message is a custom pretty error message. Might be different from the original error message
	Message string `json:"message"`
	// Details are some details to extend the error information
	Details map[string]interface{} `json:"details"`
	// Code is a constant error code to let consumers referring to the rror by its code.
	// The message might be changed, but the code should not
	Code string `json:"code"`
	// Suggestion is the advise to users how to fix the error.
	Suggestion string `json:"suggestion"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) Unwrap() error {
	return e.Original
}

type Errors []Error

// NoPrettifierErr occurs when Prettify or NewStackWithStrategy cannot match a prettifier to the given error.
var NoPrettifierErr = errors.New("prettifier is not defined")

type StackOptions struct {
	Strategy Strategy
}

type StackOptionFn func(*StackOptions)

// Prettify builds a Errors from the given error and all the wrapped errors: prettify, unwrap, prettify, repeat.
// By default, PrettifierStrategy is used. You can override the strategy using the options:
// Prettify(err, func(o *StackOptions) { o.Strategy = SomeStrategy{} })
//
// The function expects the strategy returns a prettifier for each error.
// NoPrettifierErr is returned if strategy does not match a prettifier for the error or any wrapped error.
func Prettify(err error, opts ...StackOptionFn) (Errors, error) {
	options := &StackOptions{
		Strategy: PrettifierStrategy{},
	}

	for _, fn := range opts {
		fn(options)
	}

	next := err
	stack := Errors{}

	for next != nil {
		prettifier := options.Strategy.Match(next)
		if prettifier == nil {
			return nil, NoPrettifierErr
		}
		stack = append(stack, prettifier.Prettify())

		next = errors.Unwrap(next)
	}

	return stack, nil
}

// StandardErr prettifies error instance
func StandardErr(err error) error {
	stack, err := Prettify(err)
	if err != nil {
		return err
	}

	msg, err := TextFormatter(stack, DefaultTextFormat)
	if err != nil {
		return err
	}

	return errors.New("\n" + msg)
}
