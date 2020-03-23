package http

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
	"github.com/sirupsen/logrus"
)

// NewCaller constructs a new HTTP caller
func NewCaller() *Caller {
	return &Caller{
		ctx: context.Background(),
	}
}

// Caller represents the caller constructor
type Caller struct {
	ctx context.Context
}

// Context sets the given context as the active management context
func (caller *Caller) Context(ctx context.Context) {
	caller.ctx = ctx
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return "http"
}

// Dial constructs a new caller for the given host
func (caller *Caller) Dial(schema schema.Service, functions specs.CustomDefinedFunctions, opts schema.Options) (protocol.Call, error) {
	logger := logger.FromCtx(caller.ctx, logger.Protocol)
	logger.WithFields(logrus.Fields{
		"service": schema.GetName(),
		"host":    schema.GetHost(),
	}).Info("Constructing new HTTP caller")

	options, err := ParseCallerOptions(opts)
	if err != nil {
		return nil, err
	}

	methods := make(map[string]*Method, len(schema.GetMethods()))

	for _, method := range schema.GetMethods() {
		request, endpoint, err := GetMethodEndpoint(method)
		if err != nil {
			return nil, err
		}

		references, err := TemplateReferences(endpoint, functions)
		if err != nil {
			return nil, err
		}

		methods[method.GetName()] = &Method{
			name:       method.GetName(),
			request:    request,
			endpoint:   endpoint,
			references: references,
		}
	}

	result := &Call{
		ctx:     caller.ctx,
		logger:  logger,
		service: schema.GetName(),
		host:    schema.GetHost(),
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
	ctx     context.Context
	logger  *logrus.Logger
	service string
	host    string
	methods map[string]*Method
	proxy   *httputil.ReverseProxy
}

// GetMethods returns the available methods within the HTTP caller
func (call *Call) GetMethods() []protocol.Method {
	result := make([]protocol.Method, 0, len(call.methods))

	for _, method := range call.methods {
		result = append(result, method)
	}

	return result
}

// GetMethod attempts to return a method matching the given name
func (call *Call) GetMethod(name string) protocol.Method {
	for _, method := range call.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// SendMsg calls the configured host and attempts to call the given endpoint with the given headers and stream
func (call *Call) SendMsg(ctx context.Context, rw protocol.ResponseWriter, pr *protocol.Request, refs *refs.Store) error {
	request := http.MethodGet
	url, err := url.Parse(call.host)
	if err != nil {
		return err
	}

	if pr.Method != nil {
		method := call.methods[pr.Method.GetName()]
		if method == nil {
			return trace.New(trace.WithMessage("unkown method '%s' for service '%s'", pr.Method, call.service))
		}

		endpoint := LookupEndpointReferences(method, refs)
		if endpoint != "" {
			url.Path = endpoint
		}

		request = method.request
	}

	call.logger.WithFields(logrus.Fields{
		"url":     url,
		"service": call.service,
		"method":  request,
	}).Debug("Calling HTTP caller")

	req, err := http.NewRequestWithContext(ctx, request, url.String(), pr.Body)
	if err != nil {
		return err
	}

	req.Header = CopyMetadataHeader(pr.Header)
	res := NewProtocolResponseWriter(ctx, rw)

	call.proxy.ServeHTTP(res, req)
	rw.Header().Append(CopyHTTPHeader(res.Header()))

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	call.logger.WithField("host", call.host).Info("Closing HTTP caller")
	return nil
}

// LookupEndpointReferences looks up the references within the given endpoint and returns the newly constructed endpoint
func LookupEndpointReferences(method *Method, store *refs.Store) string {
	result := method.endpoint

	for _, prop := range method.references {
		ref := store.Load(prop.Reference.Resource, prop.Reference.Path)
		if ref == nil || prop.Type != types.TypeString {
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
func TemplateReferences(value string, functions specs.CustomDefinedFunctions) ([]*specs.Property, error) {
	references := ReferenceLookup.FindAllString(value, -1)
	result := make([]*specs.Property, 0, len(references))
	for _, key := range references {
		path := key[1:]
		property := &specs.Property{
			Path: key,
			Reference: &specs.PropertyReference{
				Resource: ".request",
				Path:     path,
			},
		}

		result = append(result, property)
	}

	return result, nil
}

// GetMethodEndpoint attempts to find the endpoint for the given method.
// Empty values are returned when a empty method name is given.
func GetMethodEndpoint(method schema.Method) (string, string, error) {
	options := method.GetOptions()

	request := options[MethodOption]
	endpoint := options[EndpointOption]

	return request, endpoint, nil
}
