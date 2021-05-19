package manager

import (
	"github.com/jexia/semaphore/v2"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

// FlowOptions represents a collection of options which are used
// during the construction of a flow manager.
type FlowOptions struct {
	semaphore.Options
	stack       functions.Collection
	services    specs.ServiceList
	discoveries specs.ServiceDiscoveryClients
}

// FlowOption applies the given options to the apply options object.
type FlowOption func(*FlowOptions)

// WithFlowFunctions sets the given functions
func WithFlowFunctions(stack functions.Collection) FlowOption {
	return func(options *FlowOptions) {
		options.stack = stack
	}
}

// WithFlowServices sets the given services
func WithFlowServices(services specs.ServiceList) FlowOption {
	return func(options *FlowOptions) {
		options.services = services
	}
}

// WithServiceDiscoveries sets the given discoveries
func WithServiceDiscoveries(discoveries specs.ServiceDiscoveryClients) FlowOption {
	return func(options *FlowOptions) {
		options.discoveries = discoveries
	}
}

// WithFlowOptions sets the given options
func WithFlowOptions(conf semaphore.Options) FlowOption {
	return func(options *FlowOptions) {
		options.Options = conf
	}
}

// NewFlowOptions constructs a new endpoint option object from the passed options
func NewFlowOptions(opts ...FlowOption) FlowOptions {
	result := FlowOptions{}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}
