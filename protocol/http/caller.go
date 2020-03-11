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
func (caller *Caller) New(schema schema.Service, serviceMethod string, opts schema.Options) (protocol.Call, error) {
	log.WithFields(log.Fields{
		"service": schema.GetName(),
		"method":  serviceMethod,
	}).Info("Constructing new HTTP caller")

	options, err := ParseCallerOptions(opts)
	if err != nil {
		return nil, err
	}

	_, err = url.Parse(schema.GetHost())
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

	return &Call{
		service:  schema.GetName(),
		host:     schema.GetHost(),
		method:   method,
		endpoint: endpoint,
		proxy:    NewProxy(options),
	}, nil
}

// Call represents the HTTP caller implementation
type Call struct {
	service  string
	host     string
	method   string
	endpoint string
	proxy    *httputil.ReverseProxy
}

// Call opens a new connection to the configured host and attempts to send the given headers and stream
func (call *Call) Call(rw protocol.ResponseWriter, incoming *protocol.Request, refs *refs.Store) error {
	url, err := url.Parse(call.host)
	if err != nil {
		return err
	}

	endpoint := call.endpoint

	// FIXME: this is a prototype which has to be replaces with a permanent system
	for _, key := range ReferenceLookup.FindAllString(endpoint, -1) {
		path := key[1:]
		ref := refs.Load(specs.InputResource, path)
		val := ""

		if ref != nil {
			str, is := ref.Value.(string)
			if is {
				val = str
			}
		}

		endpoint = strings.Replace(endpoint, string(key), val, 1)
	}

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
