package http

import (
	"net/http"
	"strings"

	"github.com/jexia/maestro/protocol"
)

// CopyProtocolHeader copies the given protocol header into a HTTP header
func CopyProtocolHeader(header protocol.Header) http.Header {
	result := http.Header{}
	for key, val := range header {
		result.Set(key, val)
	}

	return result
}

// CopyHTTPHeader copies the given HTTP header into a protocol header
func CopyHTTPHeader(header http.Header) protocol.Header {
	result := protocol.Header{}
	for key, val := range header {
		result.Set(key, strings.Join(val, ";"))
	}

	return result
}

// NewProtocolResponseWriter constructs a new HTTP response writer of the given protocol response writer
func NewProtocolResponseWriter(rw protocol.ResponseWriter) *ProtocolResponseWriter {
	return &ProtocolResponseWriter{
		header:   CopyProtocolHeader(rw.Header()),
		protocol: rw,
	}
}

// A ProtocolResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ProtocolResponseWriter struct {
	header   http.Header
	protocol protocol.ResponseWriter
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (rw *ProtocolResponseWriter) Header() http.Header {
	return rw.header
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *ProtocolResponseWriter) Write(bb []byte) (int, error) {
	return rw.protocol.Write(bb)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (rw *ProtocolResponseWriter) WriteHeader(status int) {
	rw.protocol.WriteHeader(status)
}

// NewRequest constructs a new protocol request of the given http request
func NewRequest(req *http.Request) *protocol.Request {
	return &protocol.Request{
		Context: req.Context(),
		Header:  protocol.Header{},
		Body:    req.Body,
	}
}

// NewResponseWriter constructs a new HTTP response writer of the given HTTP response writer
func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		header: make(protocol.Header),
		writer: rw,
	}
}

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ResponseWriter struct {
	header protocol.Header
	writer http.ResponseWriter
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (rw *ResponseWriter) Header() protocol.Header {
	return rw.header
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *ResponseWriter) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (rw *ResponseWriter) WriteHeader(status int) {
	rw.writer.WriteHeader(status)
}
