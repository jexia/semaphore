package flow

import (
	"bytes"
	"context"

	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/sirupsen/logrus"
)

// NewRequest constructs a new request for the given codec and header manager
func NewRequest(functions functions.Stack, codec codec.Manager, metadata *metadata.Manager) *Request {
	return &Request{
		functions: functions,
		codec:     codec,
		metadata:  metadata,
	}
}

// CallOptions represents the available options that could be used to construct a new flow caller
type CallOptions struct {
	Transport transport.Call
	Method    transport.Method
	Request   *Request
	Response  *Request
}

// NewCall constructs a new flow caller from the given transport caller and
func NewCall(ctx instance.Context, node *specs.Node, options *CallOptions) Call {
	return &Caller{
		ctx:       ctx,
		node:      node,
		transport: options.Transport,
		method:    options.Method,
		request:   options.Request,
		response:  options.Response,
	}
}

// Request represents a codec and header manager
type Request struct {
	functions functions.Stack
	codec     codec.Manager
	metadata  *metadata.Manager
}

// Caller represents a flow transport caller
type Caller struct {
	ctx        instance.Context
	node       *specs.Node
	method     transport.Method
	transport  transport.Call
	references []*specs.Property
	request    *Request
	response   *Request
}

// References returns the references inside the configured transport caller
func (caller *Caller) References() []*specs.Property {
	return caller.references
}

// Do is called by the flow manager to call the configured service
func (caller *Caller) Do(ctx context.Context, store refs.Store) error {
	bb := bytes.NewBuffer(make([]byte, 0))
	w := transport.NewResponseWriter(bb)
	r := &transport.Request{
		Method: caller.method,
	}

	if caller.request != nil {
		if caller.request.functions != nil {
			err := ExecuteFunctions(caller.request.functions, store)
			if err != nil {
				return err
			}
		}

		if caller.request.metadata != nil {
			r.Header = caller.request.metadata.Marshal(store)
		}

		if caller.request.codec != nil {
			body, err := caller.request.codec.Marshal(store)
			if err != nil {
				return err
			}

			r.Body = body
		}
	}

	result := make(chan error, 1)
	defer close(result)

	go func() {
		if caller.transport == nil {
			result <- nil
			return
		}

		result <- caller.transport.SendMsg(ctx, w, r, store)
	}()

	if caller.response != nil {
		if caller.response.codec != nil {
			err := caller.response.codec.Unmarshal(bb, store)
			if err != nil {
				return err
			}
		}
	}

	err := <-result
	if err != nil {
		caller.ctx.Logger(logger.Flow).WithFields(logrus.Fields{
			"node": caller.node.Name,
			"err":  err,
		}).Error("Service error")

		return err
	}

	if caller.response != nil {
		if caller.response.functions != nil {
			err := ExecuteFunctions(caller.response.functions, store)
			if err != nil {
				return err
			}
		}

		if caller.response.metadata != nil {
			caller.response.metadata.Unmarshal(w.Header(), store)
		}
	}

	return nil
}

// ExecuteFunctions executes the given functions and writes the results to the given store
func ExecuteFunctions(stack functions.Stack, store refs.Store) error {
	for key, function := range stack {
		resource := template.JoinPath(template.StackResource, key)
		err := function.Fn(refs.NewPrefixStore(store, resource, ""))
		if err != nil {
			return err
		}
	}

	return nil
}
