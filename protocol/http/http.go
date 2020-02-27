package http

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/julienschmidt/httprouter"
)

// NewCaller constructs a new caller for the given host
func NewCaller(host string, opts specs.Options) (protocol.Caller, error) {
	options, err := ParseCallerOptions(opts)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	url.Path = options.Endpoint

	return &Caller{
		method: options.Method,
		url:    url.String(),
		proxy: &httputil.ReverseProxy{
			Director: func(*http.Request) {},
		},
	}, nil
}

// Caller represents the HTTP caller implementation
type Caller struct {
	method string
	url    string
	proxy  *httputil.ReverseProxy
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return "http"
}

// Call opens a new connection to the configured host and attempts to send the given headers and stream
func (caller *Caller) Call(rw protocol.ResponseWriter, incoming protocol.Request, refs *refs.Store) error {
	req, err := http.NewRequestWithContext(incoming.Context, caller.method, caller.url, incoming.Body)
	if err != nil {
		return err
	}

	req.Header = CopyProtocolHeader(incoming.Header)
	caller.proxy.ServeHTTP(NewProtocolResponseWriter(rw), req)

	return nil
}

// Close closes the given caller
func (caller *Caller) Close() error {
	return nil
}

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) (protocol.Listener, error) {
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
	listener.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listener.mutex.RLock()
		listener.router.ServeHTTP(w, r)
		listener.mutex.RUnlock()
	})

	return listener.server.ListenAndServe()
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(endpoints []*protocol.Endpoint) error {
	router := httprouter.New()

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return err
		}

		router.Handle(options.Method, options.Endpoint, Handle(endpoint))
	}

	listener.mutex.Lock()
	listener.router = router
	listener.mutex.Unlock()

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	return listener.server.Close()
}

// Handle constructs a new handle function for the given endpoint to the given flow
func Handle(endpoint *protocol.Endpoint) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		refs := endpoint.Flow.NewStore()
		err := endpoint.Request.Unmarshal(r.Body, refs)
		if err != nil {
			// TODO: handle err
		}

		defer r.Body.Close()
		err = endpoint.Flow.Call(r.Context(), refs)
		if err != nil {
			// TODO: handle err
			return
		}

		reader, err := endpoint.Response.Marshal(refs)
		if err != nil {
			// TODO: handle err
			return
		}

		_, err = io.Copy(w, reader)
		if err != nil {
			// TODO: handle err
			return
		}
	}
}
