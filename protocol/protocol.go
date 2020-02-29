package protocol

import (
	"context"
	"io"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
)

// Endpoint represents a protocol listener endpoint
type Endpoint struct {
	Listener string
	Flow     *flow.Manager
	Request  codec.Manager
	Response codec.Manager
	Forward  Call
	Options  specs.Options
}

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
	Method  string
	Header  Header
	Body    io.Reader
	Context context.Context
}

// Caller constructs new calls which could be used to call services
type Caller interface {
	Name() string
	New(host string, schema schema.Service, options specs.Options) (Call, error)
}

// Call is a preconfigured interface for a single service
type Call interface {
	Call(writer ResponseWriter, request *Request, refs *refs.Store) error
	Close() error
}

// Listener specifies the listener implementation
type Listener interface {
	Name() string
	Serve() error
	Close() error
	Handle([]*Endpoint) error
}
