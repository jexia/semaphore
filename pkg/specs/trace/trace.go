package trace

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// Options able to be passed when constructing a tracing error
type Options struct {
	expr    hcl.Expression
	message string
}

// Option definition
type Option func(*Options)

// New returns a stack trace for the given parameter
func New(opts ...Option) error {
	options := Options{}
	for _, option := range opts {
		option(&options)
	}

	if options.expr == nil {
		return errors.New(options.message)
	}

	r := options.expr.Range()
	position := fmt.Sprintf("%s:%d", r.Filename, r.Start.Line)

	return fmt.Errorf("%s %s", position, options.message)
}

// WithExpression sets the given expression as a trace option
func WithExpression(expr hcl.Expression) Option {
	return func(options *Options) {
		options.expr = expr
	}
}

// WithMessage sets the given property formatted message
func WithMessage(format string, params ...interface{}) Option {
	return func(options *Options) {
		options.message = fmt.Sprintf(format, params...)
	}
}
