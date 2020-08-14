package http

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
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
func (caller *Caller) Dial(service *specs.Service, functions functions.Custom, opts specs.Options) (transport.Call, error) {
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
		ctx:     caller.ctx,
		service: service.Name,
		host:    service.Host,
		proxy:   NewProxy(options),
		methods: methods,
	}

	return result, nil
}

// Method represents a service method
type Method struct {
	name       string
	request    string
	endpoint   string
	references []*specs.Property
}

// GetName returns the method name
func (method *Method) GetName() string {
	return method.name
}

// References returns the available method references
func (method *Method) References() []*specs.Property {
	if method.references == nil {
		return make([]*specs.Property, 0)
	}

	return method.references
}

// Call represents the HTTP caller implementation
type Call struct {
	ctx     *broker.Context
	service string
	host    string
	methods map[string]*Method
	proxy   *httputil.ReverseProxy
}

// GetMethods returns the available methods within the HTTP caller
func (call *Call) GetMethods() []transport.Method {
	result := make([]transport.Method, 0, len(call.methods))

	for _, method := range call.methods {
		result = append(result, method)
	}

	return result
}

// GetMethod attempts to return a method matching the given name
func (call *Call) GetMethod(name string) transport.Method {
	for _, method := range call.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// SendMsg calls the configured host and attempts to call the given endpoint with the given headers and stream
func (call *Call) SendMsg(ctx context.Context, rw transport.ResponseWriter, pr *transport.Request, refs references.Store) error {
	request := http.MethodGet
	url, err := url.Parse(call.host)
	if err != nil {
		return err
	}

	if pr.Method != nil {
		method := call.methods[pr.Method.GetName()]
		if method == nil {
			return trace.New(trace.WithMessage("unknown method '%s' for service '%s'", pr.Method, call.service))
		}

		endpoint := LookupEndpointReferences(method, refs)
		if endpoint != "" {
			url.Path = endpoint
		}

		request = method.request
	}

	logger.Debug(call.ctx, "calling HTTP caller",
		zap.String("url", url.String()),
		zap.String("method", request),
	)

	req, err := http.NewRequestWithContext(ctx, request, url.String(), pr.Body)
	if err != nil {
		return err
	}

	req.Header = CopyMetadataHeader(pr.Header)

	if ContentTypes[pr.Codec] != "" {
		req.Header.Add("Content-Type", ContentTypes[pr.Codec])
	}

	res := NewTransportResponseWriter(ctx, rw)

	go func() {
		defer rw.Close()
		call.proxy.ServeHTTP(res, req)
	}()

	res.AwaitStatus()

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	logger.Info(call.ctx, "closing HTTP caller", zap.String("host", call.host))
	return nil
}

// LookupEndpointReferences looks up the references within the given endpoint and returns the newly constructed endpoint
func LookupEndpointReferences(method *Method, store references.Store) string {
	result := method.endpoint

	for _, prop := range method.references {
		ref := store.Load(prop.Reference.Resource, prop.Reference.Path)
		if ref == nil || prop.Type != types.String {
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
			Path: key,
			Reference: &specs.PropertyReference{
				Resource: ".params",
				Path:     path,
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
