package http

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/transport"
)

// AppendHTTPHeader appends the given HTTP header into a transport header
func AppendHTTPHeader(dest metadata.MD, src http.Header) {
	for key, vals := range src {
		dest[strings.ToLower(key)] = strings.Join(vals, ";")
	}
}

// CopyHTTPHeader copies the given HTTP header into a transport header
func CopyHTTPHeader(source http.Header) metadata.MD {
	result := metadata.MD{}
	for key, vals := range source {
		result[strings.ToLower(key)] = strings.Join(vals, ";")
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

var namedParameters = regexp.MustCompile(`[\:\*]([a-zA-Z\d\^\&\%\$@\_\-\.]+)`)

// NamedParameters returns the available named parameters inside the given url
func NamedParameters(url string) []string {
	matched := namedParameters.FindAllStringSubmatch(url, -1)
	result := make([]string, len(matched))

	for index, key := range matched {
		result[index] = key[1]
	}

	return result
}

// RawNamedParameters returns the available named parameters including the selector
func RawNamedParameters(url string) []string {
	return namedParameters.FindAllString(url, -1)
}
