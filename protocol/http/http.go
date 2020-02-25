package http

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
)

// NewCaller constructs a new caller for the given host
func NewCaller(url string, options specs.Options) protocol.Caller {
	return &Caller{
		url: url,
		proxy: &httputil.ReverseProxy{
			Director: func(*http.Request) {},
		},
	}
}

// Caller represents the HTTP caller implementation
type Caller struct {
	method string
	url    string
	proxy  *httputil.ReverseProxy
}

// Call opens a new connection to the configured host and attempts to send the given headers and stream
func (caller *Caller) Call(rw protocol.ResponseWriter, incoming protocol.Request) error {
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
func NewListener(addr string, options specs.Options) protocol.Listener {
	return &Listener{
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

// Listener represents a HTTP listener
type Listener struct {
	server *http.Server
}

// Serve opens the HTTP listener and calls the given handler function on reach request
func (listener *Listener) Serve(handler protocol.Handler) error {
	listener.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := protocol.Request{
			Context: r.Context(),
			Header:  CopyHTTPHeader(w.Header()),
			Body:    r.Body,
		}

		handler(NewResponseWriter(w), req)
	})

	return listener.server.ListenAndServe()
}

// Close closes the given listener
func (listener *Listener) Close() error {
	return listener.server.Close()
}
