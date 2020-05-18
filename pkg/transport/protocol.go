package transport

import (
	"context"
	"io"

	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

// ResponseWriter specifies the response writer implementation which could be used to both proxy forward a request or used to call a service
type ResponseWriter interface {
	Header() metadata.MD
	Write([]byte) (int, error)
}

// Request represents the request object given to a caller implementation used to make calls
type Request struct {
	Codec  string
	Header metadata.MD
	Method Method
	Body   io.Reader
}

// Callers represents a collection of callers
type Callers []Caller

// Get attempts to return a caller with the given name
func (collection Callers) Get(name string) Caller {
	for _, caller := range collection {
		if caller.Name() == name {
			return caller
		}
	}

	return nil
}

// NewCaller constructs a new caller with the given context
type NewCaller func(ctx instance.Context) Caller

// Caller constructs new calls which could be used to call services
type Caller interface {
	Name() string
	Dial(service *specs.Service, functions functions.Custom, options specs.Options) (Call, error)
}

// Call is a preconfigured interface for a single service
type Call interface {
	SendMsg(ctx context.Context, writer ResponseWriter, request *Request, refs refs.Store) error
	GetMethods() []Method
	GetMethod(name string) Method
	Close() error
}

// Method represents a call method which could be called
type Method interface {
	GetName() string
	References() []*specs.Property
}

// Listeners represents a collection of listeners
type Listeners []Listener

// Get attempts to return a listener with the given name
func (collection Listeners) Get(name string) Listener {
	for _, listener := range collection {
		if listener.Name() == name {
			return listener
		}
	}

	return nil
}

// Flow represents a flow which could be called by a transport
type Flow interface {
	NewStore() refs.Store
	GetName() string
	Do(ctx context.Context, refs refs.Store) error
	Wait()
}

// Forward represents the forward specifications
type Forward struct {
	Header  specs.Header
	Service *specs.Service
}

// Endpoint represents a transport listener endpoint
type Endpoint struct {
	Listener string
	Flow     Flow
	Request  *specs.ParameterMap
	Response *specs.ParameterMap
	Forward  *Forward
	Options  specs.Options
}

// NewListener constructs a new listener with the given context
type NewListener func(ctx instance.Context) Listener

// Listener specifies the listener implementation
type Listener interface {
	Name() string
	Serve() error
	Close() error
	Handle([]*Endpoint, map[string]codec.Constructor) error
}
