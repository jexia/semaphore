package endpoints

import (
	"fmt"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/manager"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// NewOptions constructs a new endpoint option object from the passed options
func NewOptions(opts ...EndpointOption) Options {
	result := Options{}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(&result)
	}

	return result
}

// Options represents a collection of options which are used
// construct and endpoints.
type Options struct {
	semaphore.Options
	stack    functions.Collection
	services specs.ServiceList
}

// EndpointOption applies the given options to the apply options object.
type EndpointOption func(*Options)

// WithFunctions sets the given functions
func WithFunctions(stack functions.Collection) EndpointOption {
	return func(options *Options) {
		options.stack = stack
	}
}

// WithServices sets the given services
func WithServices(services specs.ServiceList) EndpointOption {
	return func(options *Options) {
		options.services = services
	}
}

// WithCore sets the given core options
func WithCore(conf semaphore.Options) EndpointOption {
	return func(options *Options) {
		options.Options = conf
	}
}

// Transporters constructs a new transport Endpoints list from the given endpoints and options
func Transporters(ctx *broker.Context, endpoints specs.EndpointList, flows specs.FlowListInterface, opts ...EndpointOption) (transport.EndpointList, error) {
	if ctx == nil {
		return nil, nil
	}

	options := NewOptions(opts...)

	results := make(transport.EndpointList, len(endpoints))
	logger.Debug(ctx, "constructing endpoints")

	for index, endpoint := range endpoints {
		logger.Debug(ctx, "constructing flow manager", zap.String("flow", endpoint.Flow))

		selected := flows.Get(endpoint.Flow)
		manager, err := manager.NewFlow(ctx, selected,
			manager.WithFlowFunctions(options.stack),
			manager.WithFlowServices(options.services),
			manager.WithFlowOptions(options.Options),
		)

		if err != nil {
			return nil, fmt.Errorf("failed to construct flow: %w", err)
		}

		forward, err := forwarder(selected.GetForward(), options)
		if err != nil {
			return nil, fmt.Errorf("failed to construct flow caller: %w", err)
		}

		results[index] = transport.NewEndpoint(endpoint.Listener, manager, forward, endpoint.Options, selected.GetInput(), selected.GetOutput())
	}

	return results, nil
}

// newForward constructs a flow caller for the given call.
func forwarder(call *specs.Call, options Options) (*transport.Forward, error) {
	if call == nil {
		return nil, nil
	}

	service := options.services.Get(call.Service)
	if service == nil {
		return nil, ErrUnknownService{Service: call.Service}
	}

	result := &transport.Forward{
		Service: service,
	}

	if call.Request != nil {
		result.Schema = call.Request.Header
	}

	return result, nil
}
