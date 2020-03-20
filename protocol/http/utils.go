package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/protocol"
)

// CopyHTTPHeader copies the given HTTP header into a protocol header
func CopyHTTPHeader(source http.Header) metadata.MD {
	result := metadata.MD{}
	for key, vals := range source {
		result[key] = strings.Join(vals, ";")
	}

	return result
}

// SetHTTPHeader copies the given protocol header into a HTTP header
func SetHTTPHeader(writer http.Header, metadata metadata.MD) {
	for key, val := range metadata {
		writer.Set(key, val)
	}
}

// CopyMetadataHeader copies the given protocol header into a HTTP header
func CopyMetadataHeader(header metadata.MD) http.Header {
	result := http.Header{}
	for key, val := range header {
		result.Set(key, val)
	}

	return result
}

// NewProtocolResponseWriter constructs a new HTTP response writer of the given protocol response writer
func NewProtocolResponseWriter(ctx context.Context, rw protocol.ResponseWriter) *ProtocolResponseWriter {
	return &ProtocolResponseWriter{
		header:   http.Header{},
		protocol: rw,
	}
}

// A ProtocolResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ProtocolResponseWriter struct {
	header   http.Header
	protocol protocol.ResponseWriter
	status   int
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
	rw.status = status
}

// NewRequest constructs a new protocol request of the given http request
func NewRequest(req *http.Request) *protocol.Request {
	return &protocol.Request{
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
