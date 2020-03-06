package http

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

const (
	// EndpointOption represents the HTTP endpoints option key
	EndpointOption = "http.endpoint"
	// MethodOption represents the HTTP method option key
	MethodOption = "http.method"
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
func (caller *Caller) New(host string, serviceMethod string, schema schema.Service, opts specs.Options) (protocol.Call, error) {
	log.WithField("host", host).Info("Constructing new HTTP caller")

	callerOptions, err := ParseCallerOptions(opts)
	if err != nil {
		return nil, err
	}

	_, err = url.Parse(host)
	if err != nil {
		return nil, err
	}

	method := schema.GetMethod(serviceMethod)
	if method == nil {
		return nil, trace.New(trace.WithMessage("service method not found '%s'.'%s'", schema.GetName(), serviceMethod))
	}

	methodOptions := method.GetOptions()

	return &Call{
		host:     host,
		method:   methodOptions[MethodOption],
		endpoint: methodOptions[EndpointOption],
		proxy: &httputil.ReverseProxy{
			Director:      func(*http.Request) {},
			FlushInterval: callerOptions.FlushInterval,
		},
	}, nil
}

// Call represents the HTTP caller implementation
type Call struct {
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

	url.Path = call.endpoint
	log.WithField("url", url).Debug("Calling HTTP caller")

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

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) (protocol.Listener, error) {
	log.WithField("add", addr).Info("Constructing new HTTP listener")

	options, err := ParseEndpointOptions(opts)
	if err != nil {
		return nil, err
	}

	return &Listener{
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
		},
	}, nil
}

// Listener represents a HTTP listener
type Listener struct {
	server *http.Server
	mutex  sync.RWMutex
	router http.Handler
}

// Name returns the name of the given listener
func (listener *Listener) Name() string {
	return "http"
}

// Serve opens the HTTP listener and calls the given handler function on reach request
func (listener *Listener) Serve() error {
	log.WithField("addr", listener.server.Addr).Info("Opening HTTP listener")

	listener.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listener.mutex.RLock()
		listener.router.ServeHTTP(w, r)
		listener.mutex.RUnlock()
	})

	err := listener.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(endpoints []*protocol.Endpoint) error {
	log.Info("HTTP listener received new endpoints")
	router := httprouter.New()

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return err
		}

		router.Handle(options.Method, options.Endpoint, Handle(endpoint))
	}

	log.Info("Swapping HTTP router")
	listener.mutex.Lock()
	listener.router = router
	listener.mutex.Unlock()

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	log.Info("Closing HTTP listener")
	return listener.server.Close()
}

// Handle constructs a new handle function for the given endpoint to the given flow
func Handle(endpoint *protocol.Endpoint) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Debug("New incoming HTTP request")

		defer r.Body.Close()
		var err error
		refs := endpoint.Flow.NewStore()

		if endpoint.Header != nil {
			header := CopyHTTPHeader(r.Header)
			endpoint.Header.Unmarshal(header, refs)
		}

		if endpoint.Request != nil {
			err = endpoint.Request.Unmarshal(r.Body, refs)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		err = endpoint.Flow.Call(r.Context(), refs)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if endpoint.Header != nil {
			SetHTTPHeader(w.Header(), endpoint.Header.Marshal(refs))
		}

		if endpoint.Response != nil {
			reader, err := endpoint.Response.Marshal(refs)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(w, reader)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}

		if endpoint.Forward != nil {
			err := endpoint.Forward.Call(NewResponseWriter(w), NewRequest(r), refs)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		}
	}
}
