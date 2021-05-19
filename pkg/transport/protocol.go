package transport

import (
	"context"
	"io"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/codec"
	"github.com/jexia/semaphore/v2/pkg/codec/metadata"
	"github.com/jexia/semaphore/v2/pkg/discovery"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

// ResponseWriter specifies the response writer implementation which could be used to both proxy forward a request or used to call a service
type ResponseWriter interface {
	io.WriteCloser
	Header() metadata.MD
	HeaderStatus(int)
	HeaderMessage(string)
}

// Request represents the request object given to a caller implementation used to make calls
type Request struct {
	RequestCodec  string
	ResponseCodec string
	Header        metadata.MD
	Method        Method
	Body          io.Reader
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
type NewCaller func(ctx *broker.Context) Caller

// Caller constructs new calls which could be used to call services
type Caller interface {
	Name() string
	Dial(service *specs.Service, functions functions.Custom, options specs.Options, resolver discovery.Resolver) (Call, error)
}

// Call is a preconfigured interface for a single service
type Call interface {
	SendMsg(ctx context.Context, writer ResponseWriter, request *Request, refs references.Store) error
	GetMethods() []Method
	GetMethod(name string) Method
	Close() error
}

// Method represents a call method which could be called
type Method interface {
	GetName() string
	References() []*specs.Property
}

// ListenerList represents a collection of listeners
type ListenerList []Listener

// Get attempts to return a listener with the given name
func (collection ListenerList) Get(name string) Listener {
	for _, listener := range collection {
		if listener.Name() == name {
			return listener
		}
	}

	return nil
}

// Flow represents a flow which could be called by a transport
type Flow interface {
	NewStore() references.Store
	GetName() string
	Errors() []Error
	Do(ctx context.Context, refs references.Store) error
	Wait()
}

// NewListener constructs a new listener with the given context
type NewListener func(ctx *broker.Context) Listener

// Listener specifies the listener implementation
type Listener interface {
	Name() string
	Serve() error
	Close() error
	Handle(*broker.Context, []*Endpoint, map[string]codec.Constructor) error
}
