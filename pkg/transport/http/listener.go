package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, origins []string, opts specs.Options) transport.NewListener {
	return func(parent *broker.Context) transport.Listener {
		module := broker.WithModule(parent, "listener", "http")
		ctx := logger.WithLogger(logger.WithFields(module, zap.String("listener", "http")))

		options, err := ParseListenerOptions(opts)
		if err != nil {
			logger.Error(ctx, "unable to parse HTTP listener options, unexpected error", zap.Error(err))
		}

		return &Listener{
			ctx:     ctx,
			options: options,
			origins: origins,
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
	ctx     *broker.Context
	options *ListenerOptions
	server  *http.Server
	mutex   sync.RWMutex
	router  http.Handler
	origins []string
}

// Name returns the name of the given listener
func (listener *Listener) Name() string { return "http" }

// Serve opens the HTTP listener and calls the given handler function on reach request
func (listener *Listener) Serve() (err error) {
	logger.Info(listener.ctx, "serving HTTP listener", zap.String("addr", listener.server.Addr))

	listener.server.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			listener.mutex.RLock()
			defer listener.mutex.RUnlock()

			if listener.router != nil {
				listener.router.ServeHTTP(w, r)
			}
		},
	)

	if listener.options.CertFile != "" && listener.options.KeyFile != "" {
		err = listener.server.ListenAndServeTLS(listener.options.CertFile, listener.options.KeyFile)
	} else {
		err = listener.server.ListenAndServe()
	}

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(ctx *broker.Context, endpoints []*transport.Endpoint, codecs map[string]codec.Constructor) error {
	logger.Info(listener.ctx, "HTTP listener received new endpoints")

	var (
		router  = httprouter.New()
		headers = make(UniqueStringItems)
		methods = make(UniqueStringItems)
	)

	for _, endpoint := range endpoints {

		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return fmt.Errorf("endpoint %s: %s", endpoint.Flow, err)
		}

		methods.Add(options.Method)

		ctx := logger.WithFields(ctx, zap.String("endpoint", options.Endpoint), zap.String("method", options.Method))
		handle, err := NewHandle(ctx, endpoint, options, codecs)
		if err != nil {
			return err
		}

		if endpoint.Request != nil {
			if endpoint.Request.Meta != nil {
				for header := range endpoint.Request.Meta.Params {
					headers.Add(header)
				}
			}
		}

		router.Handle(options.Method, options.Endpoint, handle.HTTPFunc)
	}

	router.GlobalOPTIONS = OptionsHandler(listener.origins, headers.Get(), methods.Get())

	listener.mutex.Lock()
	listener.router = router
	listener.mutex.Unlock()

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	logger.Info(listener.ctx, "closing HTTP listener")
	return listener.server.Close()
}

// NewHandle constructs a new handle function for the given endpoint to the given flow
func NewHandle(ctx *broker.Context, endpoint *transport.Endpoint, options *EndpointOptions, constructors map[string]codec.Constructor) (*Handle, error) {
	if constructors == nil {
		constructors = make(map[string]codec.Constructor)
	}

	constructor := constructors[options.Codec]
	if constructor == nil {
		return nil, trace.New(trace.WithMessage("codec not found '%s'", options.Codec))
	}

	err := endpoint.NewCodec(ctx, constructor)
	if err != nil {
		return nil, err
	}

	handle := &Handle{
		ctx:      ctx,
		Endpoint: endpoint,
		Options:  options,
	}

	if endpoint.Forward != nil {
		url, err := url.Parse(endpoint.Forward.Service.Host)
		if err != nil {
			return nil, trace.New(trace.WithMessage("unable to parse the proxy forward host '%s'", endpoint.Forward.Service.Host))
		}

		handle.Proxy = &Proxy{
			Header: endpoint.Forward.Meta,
			Handle: httputil.NewSingleHostReverseProxy(url),
		}
	}

	return handle, nil
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
	*transport.Endpoint
	ctx     *broker.Context
	Options *EndpointOptions
	Proxy   *Proxy
}

// HTTPFunc represents a HTTP function which could be used inside a HTTP router
func (handle *Handle) HTTPFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if handle == nil {
		return
	}

	logger.Debug(handle.ctx, "incoming HTTP request")

	defer r.Body.Close()
	store := handle.Endpoint.Flow.NewStore()

	for key, value := range r.URL.Query() {
		store.StoreValue(template.InputResource, key, strings.Join(value, ""))
	}

	for _, param := range ps {
		store.StoreValue(template.InputResource, param.Key, param.Value)
	}

	if handle.Request != nil {
		if handle.Request.Meta != nil {
			handle.Request.Meta.Unmarshal(CopyHTTPHeader(r.Header), store)
		}

		if handle.Request.Codec != nil {
			err := handle.Request.Codec.Unmarshal(r.Body, store)
			if err != nil {
				logger.Error(handle.ctx, "unexpected error while unmarshalling the request body", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	err := handle.Endpoint.Flow.Do(r.Context(), store)
	if err != nil {
		object := handle.Endpoint.Errs.Get(transport.Unwrap(err))
		if object == nil {
			logger.Error(handle.ctx, "unable to lookup error manager", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(object.ResolveStatusCode(store))
		if object.Meta != nil {
			SetHTTPHeader(w.Header(), object.Meta.Marshal(store))
		}

		if object.Codec != nil {
			reader, err := object.Codec.Marshal(store)
			if err != nil {
				logger.Error(handle.ctx, "unexpected error while marshalling the response body", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(w, reader)
			if err != nil {
				logger.Error(handle.ctx, "unexpected error copying the error message body to the client", zap.Error(err))
			}
		}

		return
	}

	if handle.Response != nil {
		if handle.Response.Meta != nil {
			SetHTTPHeader(w.Header(), handle.Response.Meta.Marshal(store))
		}

		if handle.Response.Codec != nil {
			ct, has := ContentTypes[handle.Response.Codec.Name()]
			if has {
				w.Header().Set(ContentTypeHeaderKey, ct)
			}

			reader, err := handle.Response.Codec.Marshal(store)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(w, reader)
			if err != nil {
				logger.Error(handle.ctx, "unexpected error copying the message body to the client", zap.Error(err))
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
