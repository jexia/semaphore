package http

import (
	"io"
	"net/http"
	"sync"

	"github.com/jexia/maestro/codec"
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
func (listener *Listener) Handle(endpoints []*protocol.Endpoint, codecs map[string]codec.Constructor) error {
	log.Info("HTTP listener received new endpoints")
	router := httprouter.New()

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint.Options)
		if err != nil {
			return err
		}

		log.Println(options.Method)
		router.Handle(options.Method, options.Endpoint, Handle(endpoint, options, codecs))
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

// Handle constructs a new handle function for the given endpoint to the given flow
func Handle(endpoint *protocol.Endpoint, options *EndpointOptions, constructors map[string]codec.Constructor) httprouter.Handle {
	if constructors == nil {
		constructors = make(map[string]codec.Constructor)
	}

	manager := constructors[options.Codec]
	if manager == nil {
		// TODO log
		return nil
	}

	var err error
	var request codec.Manager
	var response codec.Manager
	var header *protocol.HeaderManager

	if endpoint.Request != nil {
		request, err = manager.New(specs.InputResource, endpoint.Request)
		if err != nil {
			// TODO log
			return nil
		}
	}

	if endpoint.Response != nil {
		response, err = manager.New(specs.OutputResource, endpoint.Response)
		if err != nil {
			// TODO log
			return nil
		}

		header = protocol.NewHeaderManager(specs.OutputResource, endpoint.Response)
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Debug("New incoming HTTP request")

		defer r.Body.Close()
		var err error
		store := endpoint.Flow.NewStore()

		for _, param := range ps {
			store.StoreValue(specs.InputResource, param.Key, param.Value)
		}

		if header != nil {
			header.Unmarshal(CopyHTTPHeader(r.Header), store)
		}

		if endpoint.Request != nil {
			err = request.Unmarshal(r.Body, store)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		err = endpoint.Flow.Call(r.Context(), store)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if header != nil {
			SetHTTPHeader(w.Header(), header.Marshal(store))
		}

		if response != nil {
			reader, err := response.Marshal(store)
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
			err := endpoint.Forward.Call(NewResponseWriter(w), NewRequest(r), store)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		}
	}
}
