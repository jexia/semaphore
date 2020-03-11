package http

import (
	"strconv"
	"time"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
)

const (
	// EndpointOption represents the HTTP endpoints option key
	EndpointOption = "http_endpoint"
	// MethodOption represents the HTTP method option key
	MethodOption = "http_method"
	// FlushIntervalOption represents the flush interval option key
	FlushIntervalOption = "flush_interval"
	// TimeoutOption represents the timeout option key
	TimeoutOption = "timeout"
	// KeepAliveOption represents the keep alive option key
	KeepAliveOption = "keep_alive"
	// MaxIdleConnsOption represents the max idle connections option key
	MaxIdleConnsOption = "max_idle_conns"
)

// ListenerOptions represents the available HTTP options
type ListenerOptions struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ParseListenerOptions parses the given specs options into HTTP options
func ParseListenerOptions(options specs.Options) *ListenerOptions {
	result := &ListenerOptions{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	read, has := options["read_timeout"]
	if has {
		duration, err := time.ParseDuration(read)
		if err != nil {
			// TODO: log err
		}

		result.ReadTimeout = duration
	}

	write, has := options["write_timeout"]
	if has {
		duration, err := time.ParseDuration(write)
		if err != nil {
			// TODO: log err
		}

		result.WriteTimeout = duration
	}

	return result
}

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
	Timeout       time.Duration
	KeepAlive     time.Duration
	FlushInterval time.Duration
	MaxIdleConns  int
}

// ParseCallerOptions parses the given specs options into HTTP options
func ParseCallerOptions(options schema.Options) (*CallerOptions, error) {
	result := &CallerOptions{
		Timeout:      60 * time.Second,
		KeepAlive:    60 * time.Second,
		MaxIdleConns: 100,
	}

	flush, has := options[FlushIntervalOption]
	if has {
		duration, err := time.ParseDuration(flush)
		if err != nil {
			return nil, err
		}

		result.FlushInterval = duration
	}

	timeout, has := options[TimeoutOption]
	if has {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}

		result.Timeout = duration
	}

	keep, has := options[KeepAliveOption]
	if has {
		duration, err := time.ParseDuration(keep)
		if err != nil {
			return nil, err
		}

		result.KeepAlive = duration
	}

	idle, has := options[MaxIdleConnsOption]
	if has {
		value, err := strconv.Atoi(idle)
		if err != nil {
			return nil, err
		}

		result.MaxIdleConns = value
	}

	return result, nil
}
