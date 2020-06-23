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
	"github.com/jexia/maestro/pkg/specs/template"
)

// ResponseWriter specifies the response writer implementation which could be used to both proxy forward a request or used to call a service
type ResponseWriter interface {
	io.WriteCloser
	Header() metadata.MD
	HeaderStatus(int)
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

// WrapError wraps the given error as a on error object
func WrapError(err error, handle specs.ErrorHandle) Error {
	return &wrapper{
		err:    err,
		handle: handle,
	}
}

// Error represents a wrapped error and error specs
type Error interface {
	String() string
	Error() string
	Unwrap() error
	Handle() specs.ErrorHandle
}

type wrapper struct {
	err    error
	handle specs.ErrorHandle
}

func (w *wrapper) String() string {
	if w.err == nil {
		return ""
	}

	return w.err.Error()
}

// ParameterMap returns the error parameter map
func (w *wrapper) Handle() specs.ErrorHandle {
	return w.handle
}

func (w *wrapper) Error() string {
	if w.err == nil {
		return ""
	}

	return w.err.Error()
}

// Unwrap unwraps the given error and returns the wrapped error
func (w *wrapper) Unwrap() error {
	return w.err
}

// Flow represents a flow which could be called by a transport
type Flow interface {
	NewStore() refs.Store
	GetName() string
	Errors() []Error
	Do(ctx context.Context, refs refs.Store) Error
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
	Handle(instance.Context, []*Endpoint, map[string]codec.Constructor) error
}

// NewErrCodecCollection constructs a new codec collection for the given error objects
func NewErrCodecCollection(constructor codec.Constructor, collection []Error) (*CodecCollection, error) {
	result := &CodecCollection{
		collection: make(map[*specs.ParameterMap]codec.Manager, len(collection)),
	}

	for _, handle := range collection {
		if handle.Handle().GetError() == nil {
			continue
		}

		codec, err := constructor.New(template.ErrorResource, handle.Handle().GetError())
		if err != nil {
			return nil, err
		}

		result.collection[handle.Handle().GetError()] = codec
	}

	return result, nil
}

// CodecCollection represents a collection of parameter maps and their representing codec manager
type CodecCollection struct {
	collection map[*specs.ParameterMap]codec.Manager
}

// Get attempts to return a codec manager for the given parameter map
func (collection *CodecCollection) Get(hash *specs.ParameterMap) codec.Manager {
	return collection.collection[hash]
}
