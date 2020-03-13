package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
	log "github.com/sirupsen/logrus"
)

// NewCaller constructs a new HTTP caller
func NewCaller() *Caller {
	return &Caller{}
}

// Caller represents the caller constructor
type Caller struct {
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return "http"
}

// New constructs a new caller for the given host
func (caller *Caller) New(schema schema.Service, serviceMethod string, functions specs.CustomDefinedFunctions, opts schema.Options) (protocol.Call, error) {
	log.WithFields(log.Fields{
		"service": schema.GetName(),
		"method":  serviceMethod,
	}).Info("Constructing new HTTP caller")

	options, err := ParseCallerOptions(opts)
	if err != nil {
		return nil, err
	}

	schemaMethod := schema.GetMethod(serviceMethod)
	if schemaMethod == nil {
		return nil, trace.New(trace.WithMessage("service method not found '%s'.'%s'", schema.GetName(), serviceMethod))
	}

	methodOptions := schemaMethod.GetOptions()
	method := methodOptions[MethodOption]
	endpoint := methodOptions[EndpointOption]

	log.WithFields(log.Fields{
		"method":   method,
		"host":     schema.GetHost(),
		"endpoint": endpoint,
	}).Info("Constructing new HTTP caller")

	references, err := TemplateReferences(endpoint, functions)
	if err != nil {
		return nil, err
	}

	return &Call{
		service:    schema.GetName(),
		host:       schema.GetHost(),
		method:     method,
		endpoint:   endpoint,
		proxy:      NewProxy(options),
		references: references,
	}, nil
}

// Call represents the HTTP caller implementation
type Call struct {
	service    string
	host       string
	method     string
	endpoint   string
	proxy      *httputil.ReverseProxy
	references []*specs.Property
}

// References returns the available property references within the HTTP caller
func (call *Call) References() []*specs.Property {
	return call.references
}

// Call opens a new connection to the configured host and attempts to send the given headers and stream
func (call *Call) Call(rw protocol.ResponseWriter, incoming *protocol.Request, refs *refs.Store) error {
	url, err := url.Parse(call.host)
	if err != nil {
		return err
	}

	endpoint := LookupEndpointReferences(call, refs)
	url.Path = endpoint

	log.WithFields(log.Fields{
		"url":     url,
		"service": call.service,
		"method":  call.method,
	}).Debug("Calling HTTP caller")

	req, err := http.NewRequestWithContext(incoming.Context, call.method, url.String(), incoming.Body)
	if err != nil {
		return err
	}

	req.Header = CopyProtocolHeader(incoming.Header)
	call.proxy.ServeHTTP(NewProtocolResponseWriter(rw), req)

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	log.WithField("host", call.host).Info("Closing HTTP caller")
	return nil
}

// LookupEndpointReferences looks up the references within the given endpoint and returns the newly constructed endpoint
func LookupEndpointReferences(call *Call, store *refs.Store) string {
	result := call.endpoint

	for _, prop := range call.References() {
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
