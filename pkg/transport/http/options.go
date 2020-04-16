package http

import (
	"strconv"
	"time"

	"github.com/jexia/maestro/pkg/specs"
)

const (
	// ReadTimeoutOption represents the HTTP read timeout option key
	ReadTimeoutOption = "read_timeout"
	// WriteTimeoutOption represents the HTTP write timeout option key
	WriteTimeoutOption = "write_timeout"
	// EndpointOption represents the HTTP endpoints option key
	EndpointOption = "endpoint"
	// MethodOption represents the HTTP method option key
	MethodOption = "method"
	// CodecOption represents the HTTP listener codec option key
	CodecOption = "codec"
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
func ParseListenerOptions(options specs.Options) (*ListenerOptions, error) {
	result := &ListenerOptions{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	read, has := options[ReadTimeoutOption]
	if has {
		duration, err := time.ParseDuration(read)
		if err != nil {
			return nil, err
		}

		result.ReadTimeout = duration
	}

	write, has := options[WriteTimeoutOption]
	if has {
		duration, err := time.ParseDuration(write)
		if err != nil {
			return nil, err
		}

		result.WriteTimeout = duration
	}

	return result, nil
}

// EndpointOptions represents the available HTTP options
type EndpointOptions struct {
	Method       string
	Endpoint     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Codec        string
}

// ParseEndpointOptions parses the given specs options into HTTP options
func ParseEndpointOptions(options specs.Options) (*EndpointOptions, error) {
	result := &EndpointOptions{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Codec:        "json",
	}

	result.Method = options[MethodOption]
	result.Endpoint = options[EndpointOption]

	codec, has := options[CodecOption]
	if has {
		result.Codec = codec
	}

	read, has := options[ReadTimeoutOption]
	if has {
		duration, err := time.ParseDuration(read)
		if err != nil {
			return nil, err
		}

		result.ReadTimeout = duration
	}

	write, has := options[WriteTimeoutOption]
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
func ParseCallerOptions(options specs.Options) (*CallerOptions, error) {
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
