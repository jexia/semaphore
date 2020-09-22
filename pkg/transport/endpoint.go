package transport

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// NewEndpoint constructs a new transport endpoint.
func NewEndpoint(listener string, flow Flow, forward *Forward, options specs.Options, request, response *specs.ParameterMap) *Endpoint {
	result := &Endpoint{
		Listener: listener,
		Flow:     flow,
		Forward:  forward,
		Options:  options,
		Request:  NewObject(request, nil, nil),
		Response: NewObject(response, nil, nil),
	}

	return result
}

// NewObject constructs a new data object
func NewObject(schema *specs.ParameterMap, status *specs.Property, message *specs.Property) *Object {
	return &Object{
		Definition: schema,
		StatusCode: status,
		Message:    message,
	}
}

// Object represents a data object.
type Object struct {
	Definition *specs.ParameterMap
	StatusCode *specs.Property
	Message    *specs.Property
	Codec      codec.Manager
	Meta       *metadata.Manager
}

// ResolveStatusCode attempts to resolve the defined status code.
// If no status code property has been defined or the property is not a int64.
// Is a internal server error status returned.
func (object *Object) ResolveStatusCode(store references.Store) int {
	if object.StatusCode == nil {
		return StatusInternalErr
	}

	result := object.StatusCode.Scalar.Default
	if object.StatusCode.Reference != nil {
		ref := store.Load(object.StatusCode.Reference.Resource, object.StatusCode.Reference.Path)
		if ref != nil && ref.Value != nil {
			result = ref.Value
		}
	}

	if result == nil {
		return StatusInternalErr
	}

	return int(result.(int64))
}

// ResolveMessage attempts to resolve the defined message.
// If no message property has been defined or the property is not a string.
// Is a internal server error message returned.
func (object *Object) ResolveMessage(store references.Store) string {
	if object.Message == nil {
		return StatusMessage(StatusInternalErr)
	}

	result := object.Message.Scalar.Default
	if object.Message.Reference != nil {
		ref := store.Load(object.Message.Reference.Resource, object.Message.Reference.Path)
		if ref != nil && ref.Value != nil {
			result = ref.Value
		}
	}

	if result == nil {
		return StatusMessage(StatusInternalErr)
	}

	return result.(string)
}

// NewMeta updates the current object metadata manager
func (object *Object) NewMeta(ctx *broker.Context, resource string) {
	if object == nil || object.Definition == nil {
		return
	}

	object.Meta = metadata.NewManager(ctx, resource, object.Definition.Header)
}

// NewCodec updates the given object to use the given codec.
// Errors returned while constructing a new codec manager are returned.
func (object *Object) NewCodec(ctx *broker.Context, resource string, codec codec.Constructor) error {
	if object == nil || object.Definition == nil || codec == nil {
		return nil
	}

	manager, err := codec.New(resource, object.Definition)
	if err != nil {
		return err
	}

	object.Codec = manager
	return nil
}

// Errs represents a err object collection
type Errs map[specs.ErrorHandle]*Object

// Set appends the given object to the object collection
func (collection Errs) Set(err Error, object *Object) {
	if collection == nil || object == nil {
		return
	}

	collection[err.Ptr()] = object
}

// Get attempts to retrieve the requested object from the errs collection
func (collection Errs) Get(key Error) *Object {
	if collection == nil || key == nil {
		return nil
	}

	return collection[key.Ptr()]
}

// Forward represents the forward specifications
type Forward struct {
	Schema  specs.Header
	Meta    *metadata.Manager
	Service *specs.Service
}

// NewMeta updates the current object metadata manager
func (forward *Forward) NewMeta(ctx *broker.Context, resource string) {
	if forward == nil || forward.Schema == nil {
		return
	}

	forward.Meta = metadata.NewManager(ctx, resource, forward.Schema)
}

// EndpointList represents a collection of transport endpoints
type EndpointList []*Endpoint

// Endpoint represents a transport listener endpoint
type Endpoint struct {
	Listener string
	Flow     Flow
	Forward  *Forward
	Options  specs.Options
	Request  *Object
	Response *Object
	Errs     Errs
}

// NewCodec updates the endpoint request and response codecs and metadata managers.
// If a forwarding service is set is the request codec ignored.
func (endpoint *Endpoint) NewCodec(ctx *broker.Context, request codec.Constructor, response codec.Constructor) (err error) {
	endpoint.Request.NewMeta(ctx, template.InputResource)

	if endpoint.Forward == nil && endpoint.Request != nil {
		err = endpoint.Request.NewCodec(ctx, template.InputResource, request)
		if err != nil {
			return err
		}
	}

	if endpoint.Errs == nil {
		endpoint.Errs = Errs{}
	}

	if endpoint.Flow != nil {
		for _, handle := range endpoint.Flow.Errors() {
			object := NewObject(handle.GetResponse(), handle.GetStatusCode(), handle.GetMessage())

			object.NewMeta(ctx, template.ErrorResource)
			err = object.NewCodec(ctx, template.ErrorResource, response)
			if err != nil {
				return err
			}

			endpoint.Errs.Set(handle, object)
		}
	}

	endpoint.Response.NewMeta(ctx, template.OutputResource)
	err = endpoint.Response.NewCodec(ctx, template.OutputResource, response)
	if err != nil {
		return err
	}

	endpoint.Forward.NewMeta(ctx, template.OutputResource)
	return nil
}
