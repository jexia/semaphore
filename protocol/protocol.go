package protocol

import (
	"context"
	"io"

	"github.com/jexia/maestro/specs"
)

// A Header represents the key-value pairs.
type Header map[string]string

// Clone returns a copy of h or nil if h is nil.
func (h Header) Clone() Header {
	return h
}

// Del deletes the values associated with key.
func (h Header) Del(key string) {}

// Get gets the first value associated with the given key. If there are no values associated with the key, Get returns "".
func (h Header) Get(key string) string {
	return ""
}

// Set sets the header entries associated with key to the single element value. It replaces any existing values associated with key.
func (h Header) Set(key, value string) {}

// ResponseWriter specifies the response writer implementation which could be used to both proxy forward a request or used to call a service
type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	WriteHeader(int)
}

// Request represents the request object given to a caller implementation used to make calls
type Request struct {
	Header  Header
	Body    io.ReadCloser
	Context context.Context
}

// NewCaller constructs a new caller for the given url
type NewCaller func(url string, options specs.Options) Caller

// Caller specifies the caller implementation.
type Caller interface {
	Call(writer ResponseWriter, request Request) error
	Close() error
}

// Handler represents a handler function which should be called once a request has been received
type Handler func(writer ResponseWriter, request Request)

// NewListener constructs a new listener for the given addr
type NewListener func(addr string, options specs.Options) Listener

// Listener specifies the listener implementation
type Listener interface {
	Serve(Handler) error
	Close() error
}
