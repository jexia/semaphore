package http

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) transport.NewListener {
	options, err := ParseListenerOptions(opts)
	if err != nil {
		// TODO: log err
	}

	return func(ctx instance.Context) transport.Listener {
		return &Listener{
			ctx: ctx,
			server: &http.Server{
				Addr:         addr,
				ReadTimeout:  options.ReadTimeout,
				WriteTimeout: options.WriteTimeout,
			},
		}
	}
}

// Listener represents a HTTP listener
type Listener struct {
	ctx    instance.Context
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
	listener.ctx.Logger(logger.Transport).WithField("addr", listener.server.Addr).Info("Serving HTTP listener")

	listener.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listener.mutex.RLock()
		if listener.router != nil {
			listener.router.ServeHTTP(w, r)
		}
		listener.mutex.RUnlock()
	})

	err := listener.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(endpoints []*transport.Endpoint, codecs map[string]codec.Constructor) error {
	logger := listener.ctx.Logger(logger.Transport)
	logger.Info("HTTP listener received new endpoints")

	router := httprouter.New()

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return err
		}

		handle := NewHandle(logger, endpoint, options, codecs)
		router.Handle(options.Method, options.Endpoint, handle.HTTPFunc)
	}

	listener.mutex.Lock()
	listener.router = router
	listener.mutex.Unlock()

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	listener.ctx.Logger(logger.Transport).Info("Closing HTTP listener")
	return listener.server.Close()
}

// NewHandle constructs a new handle function for the given endpoint to the given flow
func NewHandle(logger *logrus.Logger, endpoint *transport.Endpoint, options *EndpointOptions, constructors map[string]codec.Constructor) *Handle {
	if constructors == nil {
		constructors = make(map[string]codec.Constructor)
	}

	codec := constructors[options.Codec]
	if codec == nil {
		// TODO log
		return nil
	}

	handle := &Handle{
		logger:   logger,
		Endpoint: endpoint,
		Options:  options,
	}

	if endpoint.Request != nil {
		request, err := codec.New(template.InputResource, endpoint.Request)
		if err != nil {
			// TODO log
			return nil
		}

		header := metadata.NewManager(template.InputResource, endpoint.Request.Header)
		handle.Request = &Request{
			Header: header,
			Codec:  request,
		}
	}

	if endpoint.Response != nil {
		response, err := codec.New(template.OutputResource, endpoint.Response)
		if err != nil {
			// TODO log
			return nil
		}

		header := metadata.NewManager(template.OutputResource, endpoint.Response.Header)
		handle.Response = &Request{
			Header: header,
			Codec:  response,
		}
	}

	if endpoint.Forward != nil {
		url, err := url.Parse(endpoint.Forward.Service.Host)
		if err != nil {
			return nil
		}

		header := metadata.NewManager(template.OutputResource, endpoint.Forward.Header)
		handle.Proxy = &Proxy{
			Header: header,
			Handle: httputil.NewSingleHostReverseProxy(url),
		}
	}

	return handle
}

// Proxy represents a HTTP reverse proxy
type Proxy struct {
	Handle *httputil.ReverseProxy
	Header *metadata.Manager
}

// Request represents a codec manager and header manager
type Request struct {
	Codec  codec.Manager
	Header *metadata.Manager
}

// Handle holds a endpoint its options and a optional request and response
type Handle struct {
	logger   *logrus.Logger
	Endpoint *transport.Endpoint
	Options  *EndpointOptions
	Request  *Request
	Response *Request
	Proxy    *Proxy
}

// HTTPFunc represents a HTTP function which could be used inside a HTTP router
func (handle *Handle) HTTPFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if handle == nil {
		return
	}

	handle.logger.Debug("New incoming HTTP request")

	defer r.Body.Close()
	var err error
	store := handle.Endpoint.Flow.NewStore()

	for _, param := range ps {
		store.StoreValue(template.InputResource, param.Key, param.Value)
	}

	if handle.Request != nil {
		if handle.Request.Header != nil {
			handle.Request.Header.Unmarshal(CopyHTTPHeader(r.Header), store)
		}

		if handle.Request.Codec != nil {
			err = handle.Request.Codec.Unmarshal(r.Body, store)
			if err != nil {
				handle.logger.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	err = handle.Endpoint.Flow.Call(r.Context(), store)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if handle.Response != nil {
		if handle.Response.Header != nil {
			SetHTTPHeader(w.Header(), handle.Response.Header.Marshal(store))
		}

		if handle.Response.Codec != nil {
			reader, err := handle.Response.Codec.Marshal(store)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(w, reader)
			if err != nil {
				handle.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}
	}

	if handle.Endpoint.Forward != nil {
		SetHTTPHeader(r.Header, handle.Proxy.Header.Marshal(store))
		handle.Proxy.Handle.ServeHTTP(w, r)
	}
}
