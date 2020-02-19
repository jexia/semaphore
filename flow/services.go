package flow

// Services represents a collection of services
type Services map[string]*Service

// Get attempts to fetch the given service by name
func (services Services) Get(name string) *Service {
	return services[name]
} 

// Service represents a flow service
type Service struct {
	Codec Codec
	Call Call
	Rollback Call
}