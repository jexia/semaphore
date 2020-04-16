package logger

import (
	"github.com/jexia/maestro/pkg/specs/trace"
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
	// Transport represents the transport implementations
	Transport Module = "transport"
)

// Modules holds all logging modules
var Modules = []Module{
	Core,
	Flow,
	Transport,
}

// New constructs a new logger and initializes the configured modules
func New() *Logger {
	modules := make(map[Module]*logrus.Logger, len(Modules))

	for _, module := range Modules {
		modules[module] = logrus.New()
	}

	return &Logger{
		modules: modules,
	}
}

// Logger represents a logrus Logger manager
type Logger struct {
	modules map[Module]*logrus.Logger
}

// Get attempts to fetch the logger for the given module
func (logger *Logger) Get(module Module) *logrus.Logger {
	return logger.modules[module]
}

// SetLevel attempts to set the log level for the given module
func (logger *Logger) SetLevel(module Module, value string) error {
	if module == Global {
		for _, module := range Modules {
			err := logger.SetLevel(module, value)
			if err != nil {
				return err
			}
		}

		return nil
	}

	log := logger.modules[module]
	if logger == nil {
		return trace.New(trace.WithMessage("logger not found for module '%s'", module))
	}

	level, err := logrus.ParseLevel(value)
	if err != nil {
		return err
	}

	log.SetLevel(level)
	return nil
}
