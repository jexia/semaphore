package http

import (
	"time"

	"github.com/jexia/maestro/specs"
)

// EndpointOptions represents the available HTTP options
type EndpointOptions struct {
	Method       string
	Endpoint     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ParseEndpointOptions parses the given specs options into HTTP options
func ParseEndpointOptions(options specs.Options) (*EndpointOptions, error) {
	result := &EndpointOptions{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	result.Method = options["method"]
	result.Endpoint = options["endpoint"]

	read, has := options["read_timeout"]
	if has {
		duration, err := time.ParseDuration(read)
		if err != nil {
			return nil, err
		}

		result.ReadTimeout = duration
	}

	write, has := options["write_timeout"]
	if has {
		duration, err := time.ParseDuration(write)
		if err != nil {
			return nil, err
		}

		result.WriteTimeout = duration
	}

	return result, nil
}

// CallerOptions represents the available HTTP options
type CallerOptions struct {
	FlushInterval time.Duration
}

// ParseCallerOptions parses the given specs options into HTTP options
func ParseCallerOptions(options specs.Options) (*CallerOptions, error) {
	result := &CallerOptions{}

	flush, has := options["flush_interval"]
	if has {
		duration, err := time.ParseDuration(flush)
		if err != nil {
			return nil, err
		}

		result.FlushInterval = duration
	}

	return result, nil
}
