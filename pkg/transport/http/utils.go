package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/transport"
)

// CopyHTTPHeader copies the given HTTP header into a transport header
func CopyHTTPHeader(source http.Header) metadata.MD {
	result := metadata.MD{}
	for key, vals := range source {
		result[key] = strings.Join(vals, ";")
	}

	return result
}

// SetHTTPHeader copies the given transport header into a HTTP header
func SetHTTPHeader(writer http.Header, metadata metadata.MD) {
	for key, val := range metadata {
		writer.Set(key, val)
	}
}

// CopyMetadataHeader copies the given transport header into a HTTP header
func CopyMetadataHeader(header metadata.MD) http.Header {
	result := http.Header{}
	for key, val := range header {
		result.Set(key, val)
	}

	return result
}

// NewTransportResponseWriter constructs a new HTTP response writer of the given transport response writer
func NewTransportResponseWriter(ctx context.Context, rw transport.ResponseWriter) *TransportResponseWriter {
	return &TransportResponseWriter{
		header:    http.Header{},
		transport: rw,
	}
}

// A TransportResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type TransportResponseWriter struct {
	header    http.Header
	transport transport.ResponseWriter
	status    int
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (rw *TransportResponseWriter) Header() http.Header {
	return rw.header
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *TransportResponseWriter) Write(bb []byte) (int, error) {
	return rw.transport.Write(bb)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (rw *TransportResponseWriter) WriteHeader(status int) {
	rw.status = status
}

// Status returns the response writer status
func (rw *TransportResponseWriter) Status() int {
	return rw.status
}

// NewRequest constructs a new transport request of the given http request
func NewRequest(req *http.Request) *transport.Request {
	return &transport.Request{
		Body: req.Body,
	}
}

// NewResponseWriter constructs a new HTTP response writer of the given HTTP response writer
func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		header: make(metadata.MD),
		writer: rw,
	}
}

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ResponseWriter struct {
	header metadata.MD
	writer http.ResponseWriter
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (rw *ResponseWriter) Header() metadata.MD {
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
