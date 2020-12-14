package discovery

// Service represents all the information about a discovered service.
type Service struct {
	Host   string
	Port   int
	Name   string
	ID     string
	Scheme string
}

// Updates is used to retrieve discovered services.
type Updates chan []Service

// Resolver returns a resolved hostname
type Resolver interface {
	// Resolve the current hostname.
	// Returns the address and `true` in case if the service has been resolved.
	Resolve() (address string, ok bool)
}

type ResolverFunc func() (string, bool)

func (fn ResolverFunc) Resolve() (string, bool) {
	return fn()
}
