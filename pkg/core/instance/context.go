package instance

import (
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/sirupsen/logrus"
)

// Context represents the Semaphore context passed in between modules to
type Context interface {
	Logger(logger.Module) *logrus.Logger
	SetLevel(logger.Module, string) error
}

// NewContext constructs a new context
func NewContext() Context {
	return &context{
		logger: logger.New(),
	}
}

type context struct {
	logger *logger.Logger
}

func (ctx *context) Logger(module logger.Module) *logrus.Logger {
	return ctx.logger.Get(module)
}

func (ctx *context) SetLevel(module logger.Module, level string) error {
	return ctx.logger.SetLevel(module, level)
}
