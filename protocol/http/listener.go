package http

import (
	"io"
	"net/http"
	"sync"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/header"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) protocol.Listener {
	log.WithField("add", addr).Info("Constructing new HTTP listener")

	options, err := ParseListenerOptions(opts)
	if err != nil {
		// TODO: log err
	}

	return &Listener{
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
		},
	}
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
func (listener *Listener) Handle(endpoints []*protocol.Endpoint, codecs map[string]codec.Constructor) error {
	log.Info("HTTP listener received new endpoints")
	router := httprouter.New()

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return err
		}

		handle := NewHandle(endpoint, options, codecs)
		router.Handle(options.Method, options.Endpoint, handle.HTTPFunc)
	}

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

// NewHandle constructs a new handle function for the given endpoint to the given flow
func NewHandle(endpoint *protocol.Endpoint, options *EndpointOptions, constructors map[string]codec.Constructor) *Handle {
	if constructors == nil {
		constructors = make(map[string]codec.Constructor)
	}

	codec := constructors[options.Codec]
	if codec == nil {
		// TODO log
		return nil
	}

	handle := &Handle{
		Endpoint: endpoint,
		Options:  options,
	}

	if endpoint.Request != nil {
		request, err := codec.New(specs.InputResource, endpoint.Request)
		if err != nil {
			// TODO log
			return nil
		}

		header := header.NewManager(specs.InputResource, endpoint.Request)
		handle.Request = &Request{
			Header: header,
			Codec:  request,
		}
	}

	if endpoint.Response != nil {
		response, err := codec.New(specs.OutputResource, endpoint.Response)
		if err != nil {
			// TODO log
			return nil
		}

		header := header.NewManager(specs.OutputResource, endpoint.Response)
		handle.Response = &Request{
			Header: header,
			Codec:  response,
		}
	}

	return handle
}

// Request represents a codec manager and header manager
type Request struct {
	Codec  codec.Manager
	Header *header.Manager
}

// Handle holds a endpoint its options and a optional request and response
type Handle struct {
	Endpoint *protocol.Endpoint
	Options  *EndpointOptions
	Request  *Request
	Response *Request
}

// HTTPFunc represents a HTTP function which could be used inside a HTTP router
func (handle *Handle) HTTPFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Debug("New incoming HTTP request")

	defer r.Body.Close()
	var err error
	store := handle.Endpoint.Flow.NewStore()

	for _, param := range ps {
		store.StoreValue(specs.InputResource, param.Key, param.Value)
	}

	if handle.Request != nil {
		if handle.Request.Header != nil {
			handle.Request.Header.Unmarshal(CopyHTTPHeader(r.Header), store)
		}

		if handle.Request.Codec != nil {
			err = handle.Request.Codec.Unmarshal(r.Body, store)
			if err != nil {
				log.Error(err)
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
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}
	}

	if handle.Endpoint.Forward != nil {
		err := handle.Endpoint.Forward.SendMsg(NewResponseWriter(w), NewRequest(r), store)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
