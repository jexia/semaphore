package http

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
)

// ListenerOption defines a single listener option
type ListenerOption func(*ListenerOptions) error

// ListenerOptions represents the available HTTP options
type ListenerOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	certFile     string
	keyFile      string
	origins      []string
}

// DefaultListenerOptions returns default listener configuration
func DefaultListenerOptions() *ListenerOptions {
	return &ListenerOptions{
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
	}
}

// NewListenerOptions creates listener config with provided options
func NewListenerOptions(options ...ListenerOption) (*ListenerOptions, error) {
	o := DefaultListenerOptions()

	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}

	return o, nil
}

// WithReadTimeout overrides the default read timeout
func WithReadTimeout(timeout string) ListenerOption {
	return func(options *ListenerOptions) error {
		if timeout == "" {
			return nil
		}

		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return fmt.Errorf("invalid duration for read timeout: %q", timeout)
		}

		options.readTimeout = duration

		return nil
	}
}

// WithWriteTimeout overrides default write timeout
func WithWriteTimeout(timeout string) ListenerOption {
	return func(options *ListenerOptions) error {
		if timeout == "" {
			return nil
		}

		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return fmt.Errorf("invalid duration for write timeout: %q", timeout)
		}

		options.writeTimeout = duration

		return nil
	}
}

// WithKeyFile defines key file
func WithKeyFile(path string) ListenerOption {
	return func(options *ListenerOptions) error {
		if path == "" {
			return nil
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot open key file: %q", path)
		}

		if info.IsDir() {
			return fmt.Errorf("%q is not a file", path)
		}

		options.keyFile = path

		return nil
	}
}

// WithCertFile defines certificate file
func WithCertFile(path string) ListenerOption {
	return func(options *ListenerOptions) error {
		if path == "" {
			return nil
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot open certificate file: %q", path)
		}

		if info.IsDir() {
			return fmt.Errorf("%q is not a file", path)
		}

		options.certFile = path

		return nil
	}
}

// WithOrigins sets allowed origins for incoming preflight requests
func WithOrigins(list []string) ListenerOption {
	return func(options *ListenerOptions) error {
		options.origins = list

		return nil
	}
}
