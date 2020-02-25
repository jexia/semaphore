package http

// Options represents the available HTTP options
type Options struct {
	Method   string
	Endpoint string
}

// ParseEndpoint parses the given specs options into HTTP options
func ParseEndpoint(map[string]interface{}) *Options {
	return nil
}
