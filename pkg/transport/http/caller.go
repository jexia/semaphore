package http

import (
	"strings"

	"github.com/jexia/semaphore/v2/pkg/discovery"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/labels"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
	"github.com/jexia/semaphore/v2/pkg/transport"
	"go.uber.org/zap"
)

// NewCaller constructs a new HTTP caller
func NewCaller() transport.NewCaller {
	return func(ctx *broker.Context) transport.Caller {
		return &Caller{
			ctx: ctx,
		}
	}
}

// Caller represents the caller constructor
type Caller struct {
	ctx *broker.Context
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return "http"
}

// Dial constructs a new caller for the given host
func (caller *Caller) Dial(service *specs.Service, functions functions.Custom, opts specs.Options, resolver discovery.Resolver) (transport.Call, error) {
	module := broker.WithModule(caller.ctx, "caller", "http")
	ctx := logger.WithFields(logger.WithLogger(module), zap.String("service", service.Name))

	logger.Info(ctx, "constructing new HTTP caller", zap.String("host", service.Host))

	options, err := ParseCallerOptions(opts)
	if err != nil {
		return nil, err
	}

	methods := make(map[string]*Method, len(service.Methods))

	for _, method := range service.Methods {
		request, endpoint, err := GetMethodEndpoint(method)
		if err != nil {
			return nil, err
		}

		references, err := TemplateReferences(endpoint, functions)
		if err != nil {
			return nil, err
		}

		methods[method.Name] = &Method{
			name:       method.Name,
			request:    request,
			endpoint:   endpoint,
			references: references,
		}
	}

	result := &Call{
		ctx:      caller.ctx,
		service:  service.Name,
		host:     service.Host,
		proxy:    NewProxy(options),
		methods:  methods,
		resolver: resolver,
	}

	return result, nil
}

// LookupEndpointReferences looks up the references within the given endpoint and returns the newly constructed endpoint
func LookupEndpointReferences(method *Method, store references.Store) string {
	result := method.endpoint

	for _, prop := range method.references {
		ref := store.Load(prop.Reference.String())
		if ref == nil || prop.Scalar.Type != types.String {
			result = strings.Replace(result, prop.Path, "", 1)
			continue
		}

		str, is := ref.Value.(string)
		if !is {
			result = strings.Replace(result, prop.Path, "", 1)
			continue
		}

		result = strings.Replace(result, prop.Path, str, 1)
	}

	return result
}

// TemplateReferences returns the property references within the given value
func TemplateReferences(value string, functions functions.Custom) ([]*specs.Property, error) {
	references := RawNamedParameters(value)
	result := make([]*specs.Property, 0, len(references))
	for _, key := range references {
		path := key[1:]
		property := &specs.Property{
			Path:  key,
			Label: labels.Optional,
			Template: specs.Template{
				Reference: &specs.PropertyReference{
					Resource: ".params",
					Path:     path,
				},
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		}

		result = append(result, property)
	}

	return result, nil
}

// GetMethodEndpoint attempts to find the endpoint for the given method.
// Empty values are returned when a empty method name is given.
func GetMethodEndpoint(method *specs.Method) (string, string, error) {
	options := method.Options

	request := options[MethodOption]
	endpoint := options[EndpointOption]

	return request, endpoint, nil
}
