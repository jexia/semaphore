package endpoints

import (
	"fmt"

	"github.com/jexia/semaphore/v2"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/broker/manager"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/transport"
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
	stack       functions.Collection
	services    specs.ServiceList
	discoveries specs.ServiceDiscoveryClients
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

// WithServiceDiscoveries sets the given discovery clients
func WithServiceDiscoveries(discoveries specs.ServiceDiscoveryClients) EndpointOption {
	return func(options *Options) {
		options.discoveries = discoveries
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
			manager.WithServiceDiscoveries(options.discoveries),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to construct flow: %w", err)
		}

		forward, err := forwarder(selected, options)
		if err != nil {
			return nil, fmt.Errorf("failed to construct flow caller: %w", err)
		}

		results[index] = transport.NewEndpoint(endpoint.Listener, manager, forward, endpoint.Options, selected.GetInput(), selected.GetOutput())
	}

	return results, nil
}

// newForward constructs a flow caller for the given call.
func forwarder(flow specs.FlowInterface, options Options) (*transport.Forward, error) {
	call := flow.GetForward()

	if call == nil {
		return nil, nil
	}

	service := options.services.Get(call.Service)
	if service == nil {
		return nil, ErrUnknownService{Service: call.Service}
	}

	rewrite := make([]transport.Rewrite, len(flow.GetRewrite()), len(flow.GetRewrite()))
	for index, item := range flow.GetRewrite() {
		rewriteFunc, err := transport.NewRewrite(item.Pattern, item.Template)
		if err != nil {
			return nil, err
		}

		rewrite[index] = rewriteFunc
	}

	result := &transport.Forward{
		Service: service,
		Rewrite: rewrite,
	}

	if call.Request != nil {
		result.Schema = call.Request.Header
	}

	return result, nil
}
