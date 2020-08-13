package http

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewListenerOptions(t *testing.T) {
	type test struct {
		options  []ListenerOption
		error    error
		expected *ListenerOptions
	}

	tests := map[string]test{
		"check if error is returnes when unable to parse read timeout": {
			options: []ListenerOption{
				WithReadTimeout("one second"),
			},
			error: fmt.Errorf(`invalid duration for read timeout: "one second"`),
		},
		"check if error is returnes when unable to parse write timeout": {
			options: []ListenerOption{
				WithWriteTimeout("one minute"),
			},
			error: fmt.Errorf(`invalid duration for write timeout: "one minute"`),
		},
		"check if error is returned when path to key file is invalid": {
			options: []ListenerOption{
				WithKeyFile("/invalid/path"),
			},
			error: fmt.Errorf(`cannot open key file: "/invalid/path"`),
		},
		"check if error is returned when path to certificate file is invalid": {
			options: []ListenerOption{
				WithCertFile("/invalid/path"),
			},
			error: fmt.Errorf(`cannot open certificate file: "/invalid/path"`),
		},
		"check if empty cert/key path or empty timeouts are ignored": {
			options: []ListenerOption{
				WithKeyFile(""),
				WithCertFile(""),
				WithReadTimeout(""),
				WithWriteTimeout(""),
			},
			expected: &ListenerOptions{
				readTimeout:  defaultReadTimeout,
				writeTimeout: defaultWriteTimeout,
			},
		},
		"check if default timeouts can be overriden": {
			options: []ListenerOption{
				WithReadTimeout("1s"),
				WithWriteTimeout("1m"),
			},
			expected: &ListenerOptions{
				readTimeout:  time.Second,
				writeTimeout: time.Minute,
			},
		},
		"check if origins are set": {
			options: []ListenerOption{
				WithOrigins([]string{"test.com", "example.com"}),
			},
			expected: &ListenerOptions{
				readTimeout:  defaultReadTimeout,
				writeTimeout: defaultWriteTimeout,
				origins:      []string{"test.com", "example.com"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewListenerOptions(test.options...)

			if err != test.error {
				t.Errorf("error [%v] was expected to be [%v]", err, test.error)
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("options [%v] was expected to equal [%v]", actual, test.expected)
			}
		})
	}
}
