package http

import "time"

// EndpointOptions represents the available HTTP options
type EndpointOptions struct {
	Method       string
	Endpoint     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ParseEndpointOptions parses the given specs options into HTTP options
func ParseEndpointOptions(options map[string]interface{}) (*EndpointOptions, error) {
	result := &EndpointOptions{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	method, has := options["method"].(string)
	if has {
		result.Method = method
	}

	endpoint, has := options["endpoint"].(string)
	if has {
		result.Endpoint = endpoint
	}

	read, has := options["read_timeout"].(string)
	if has {
		duration, err := time.ParseDuration(read)
		if err != nil {
			return nil, err
		}

		result.ReadTimeout = duration
	}

	write, has := options["write_timeout"].(string)
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
	Method   string
	Endpoint string
}

// ParseCallerOptions parses the given specs options into HTTP options
func ParseCallerOptions(options map[string]interface{}) (*CallerOptions, error) {
	result := &CallerOptions{}

	method, has := options["method"].(string)
	if has {
		result.Method = method
	}

	endpoint, has := options["endpoint"].(string)
	if has {
		result.Endpoint = endpoint
	}

	return result, nil
}
