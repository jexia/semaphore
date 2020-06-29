package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/jexia/maestro/internal/codec"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/core/trace"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/julienschmidt/httprouter"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) transport.NewListener {
	return func(ctx instance.Context) transport.Listener {
		options, err := ParseListenerOptions(opts)
		if err != nil {
			ctx.Logger(logger.Transport).Warnf("unable to parse HTTP listener options, unexpected error %s", err)
		}

		return &Listener{
			ctx:     ctx,
			options: options,
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
	ctx     instance.Context
	options *ListenerOptions
	server  *http.Server
	mutex   sync.RWMutex
	router  http.Handler
	headers string
}

// Name returns the name of the given listener
func (listener *Listener) Name() string {
	return "http"
}

// Serve opens the HTTP listener and calls the given handler function on reach request
func (listener *Listener) Serve() (err error) {
	listener.ctx.Logger(logger.Transport).WithField("addr", listener.server.Addr).Info("Serving HTTP listener")

	listener.server.Handler = listener.HandleCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listener.mutex.RLock()
		if listener.router != nil {
			listener.router.ServeHTTP(w, r)
		}
		listener.mutex.RUnlock()
	}))

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
func (listener *Listener) Handle(ctx instance.Context, endpoints []*transport.Endpoint, codecs map[string]codec.Constructor) error {
	listener.ctx.Logger(logger.Transport).Info("HTTP listener received new endpoints")

	router := httprouter.New()
	headers := map[string]struct{}{}

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return fmt.Errorf("endpoint %s: %s", endpoint.Flow, err)
		}

		handle, err := NewHandle(ctx, endpoint, options, codecs)
		if err != nil {
			return err
		}

		if handle.Request != nil {
			if handle.Request.Meta != nil {
				for header := range handle.Request.Meta.Params {
					headers[header] = struct{}{}
				}
			}
		}

		router.Handle(options.Method, options.Endpoint, handle.HTTPFunc)
	}

	list := make([]string, 0, len(headers))
	for header := range headers {
		list = append(list, header)
	}

	listener.mutex.Lock()
	listener.router = router
	listener.headers = strings.Join(list, ", ")
	listener.mutex.Unlock()

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	listener.ctx.Logger(logger.Transport).Info("Closing HTTP listener")
	return listener.server.Close()
}

// HandleCors handles the defining of cors headers for the incoming HTTP request
func (listener *Listener) HandleCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		method := r.Header.Get("Access-Control-Request-Method")

		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")

		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Headers", "*")
		headers.Set("Access-Control-Allow-Methods", strings.ToUpper(method))

		if r.Method != http.MethodOptions || r.Header.Get("Access-Control-Request-Method") == "" {
			h.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// NewHandle constructs a new handle function for the given endpoint to the given flow
func NewHandle(ctx instance.Context, endpoint *transport.Endpoint, options *EndpointOptions, constructors map[string]codec.Constructor) (*Handle, error) {
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
	ctx     instance.Context
	Options *EndpointOptions
	Proxy   *Proxy
}

// HTTPFunc represents a HTTP function which could be used inside a HTTP router
func (handle *Handle) HTTPFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if handle == nil {
		return
	}

	handle.ctx.Logger(logger.Transport).Debug("New incoming HTTP request")

	defer r.Body.Close()
	store := handle.Endpoint.Flow.NewStore()

	for key, value := range r.URL.Query() {
		store.StoreValue(template.InputResource, key, value)
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
				handle.ctx.Logger(logger.Transport).Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	err := handle.Endpoint.Flow.Do(r.Context(), store)
	if err != nil {
		object := handle.Endpoint.Errs.Get(transport.Unwrap(err))
		if object == nil {
			handle.ctx.Logger(logger.Transport).Error("Unable to lookup error manager")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reader, err := object.Codec.Marshal(store)
		if err != nil {
			handle.ctx.Logger(logger.Transport).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(object.ResolveStatusCode(handle.ctx, store))
		if object.Meta != nil {
			SetHTTPHeader(w.Header(), object.Meta.Marshal(store))
		}

		_, err = io.Copy(w, reader)
		if err != nil {
			handle.ctx.Logger(logger.Transport).Error(err)
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
				handle.ctx.Logger(logger.Transport).Error(err)
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
