package logger

import (
	"context"

	"github.com/jexia/maestro/specs/trace"
	"github.com/sirupsen/logrus"
)

// Module represents a logging module
type Module string

var (
	// Global represents all modules
	Global Module = "all"
	// Core represent the internal Maestro implementations
	Core Module = "core"
	// Flow represent Maestro flow manager
	Flow Module = "flow"
	// Protocol represents the protocol implementations
	Protocol Module = "protocol"
)

// Modules holds all logging modules
var Modules = []Module{
	Core,
	Flow,
	Protocol,
}

// WithValue initialises the logger context by injecting the logging modules
func WithValue(ctx context.Context) context.Context {
	for _, module := range Modules {
		ctx = context.WithValue(ctx, module, logrus.New())
	}

	return ctx
}

// SetLevel attempts to set the log level for the given module
func SetLevel(ctx context.Context, module Module, value string) error {
	if module == Global {
		for _, module := range Modules {
			err := SetLevel(ctx, module, value)
			if err != nil {
				return err
			}
		}

		return nil
	}

	logger := FromCtx(ctx, module)
	if logger == nil {
		return trace.New(trace.WithMessage("logger not found for module '%s'", module))
	}

	level, err := logrus.ParseLevel(value)
	if err != nil {
		return err
	}

	logger.SetLevel(level)
	return nil
}

// FromCtx attempts to fetch the logger from the given context
func FromCtx(ctx context.Context, module Module) *logrus.Logger {
	value := ctx.Value(module)
	if value == nil {
		return nil
	}

	logger, is := value.(*logrus.Logger)
	if !is {
		return nil
	}

	return logger
}
