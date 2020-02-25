package http

import (
	"net/http"
	"net/http/httputil"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
)

// NewResponseWriter constructs a new HTTP response writer of the given protocol response writer
func NewResponseWriter(rw protocol.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		header:   CopyHeader(rw.Header()),
		protocol: rw,
	}
}

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ResponseWriter struct {
	header   http.Header
	protocol protocol.ResponseWriter
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (rw *ResponseWriter) Header() http.Header {
	return rw.header
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *ResponseWriter) Write(bb []byte) (int, error) {
	return rw.protocol.Write(bb)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (rw *ResponseWriter) WriteHeader(status int) {
	rw.protocol.WriteHeader(status)
}

// Caller represents the HTTP caller implementation
type Caller struct {
	method string
	url    string
	proxy  *httputil.ReverseProxy
}

// Open opens a new caller for the given host
func (caller *Caller) Open(url string, options specs.Options) protocol.Caller {
	return &Caller{
		url: url,
		proxy: &httputil.ReverseProxy{
			Director: func(*http.Request) {},
		},
	}
}

// Call opens a new connection to the configured host and attempts to send the given headers and stream
func (caller *Caller) Call(rw protocol.ResponseWriter, incoming protocol.Request) error {
	req, err := http.NewRequestWithContext(incoming.Context, caller.method, caller.url, incoming.Body)
	if err != nil {
		return err
	}

	req.Header = CopyHeader(incoming.Header)
	caller.proxy.ServeHTTP(NewResponseWriter(rw), req)

	return nil
}
