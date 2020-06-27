package flow

import (
	"context"
	"errors"
	"io"

	"github.com/jexia/maestro/internal/codec"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/sirupsen/logrus"
)

// ErrAbortFlow represents the error thrown when a flow has to be aborted
var ErrAbortFlow = errors.New("abort flow")

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
	Transport      transport.Call
	Method         transport.Method
	Request        *Request
	Response       *Request
	Err            *OnError
	ExpectedStatus int
}

// NewCall constructs a new flow caller from the given transport caller and
func NewCall(ctx instance.Context, node *specs.Node, options *CallOptions) Call {
	result := &Caller{
		ctx:            ctx,
		node:           node,
		transport:      options.Transport,
		method:         options.Method,
		request:        options.Request,
		response:       options.Response,
		ExpectedStatus: options.ExpectedStatus,
		err:            options.Err,
	}

	if result.ExpectedStatus == 0 {
		result.ExpectedStatus = transport.StatusOK
	}

	return result
}

// Request represents a codec and metadata manager
type Request struct {
	functions functions.Stack
	codec     codec.Manager
	metadata  *metadata.Manager
}

// Caller represents a flow transport caller
type Caller struct {
	ctx            instance.Context
	node           *specs.Node
	method         transport.Method
	transport      transport.Call
	references     []*specs.Property
	request        *Request
	response       *Request
	err            *OnError
	ExpectedStatus int
}

// References returns the references inside the configured transport caller
func (caller *Caller) References() []*specs.Property {
	return caller.references
}

// Do is called by the flow manager to call the configured service
func (caller *Caller) Do(ctx context.Context, store refs.Store) error {
	reader, writer := io.Pipe()
	w := transport.NewResponseWriter(writer)
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

			r.Codec = caller.request.codec.Name()
			r.Body = body
		}
	}

	if caller.transport != nil {
		err := caller.transport.SendMsg(ctx, w, r, store)
		if err != nil {
			caller.ctx.Logger(logger.Flow).WithFields(logrus.Fields{
				"node": caller.node.Name,
				"err":  err,
			}).Error("Transport returned a unexpected error")

			return err
		}
	} else {
		writer.Close()
	}

	if caller.transport != nil && w.Status() != caller.ExpectedStatus {
		caller.ctx.Logger(logger.Flow).WithFields(logrus.Fields{
			"node":   caller.node.Name,
			"status": w.Status(),
		}).Error("Service returned a unexpected status, aborting flow")

		err := caller.HandleErr(w, reader, store)
		if err != nil {
			return err
		}

		return ErrAbortFlow
	}

	if caller.response != nil {
		if caller.response.codec != nil {
			err := caller.response.codec.Unmarshal(reader, store)
			if err != nil {
				return err
			}
		}

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

// HandleErr handles a thrown service error. If a error response is defined is it decoded
func (caller *Caller) HandleErr(w *transport.Writer, reader io.Reader, store refs.Store) error {
	var status interface{}
	var message interface{}

	if caller.err != nil {
		if caller.err.message != nil {
			message = caller.err.message.Default

			if caller.err.message.Reference != nil {
				ref := store.Load(caller.err.message.Reference.Resource, caller.err.message.Reference.Path)
				message = ref.Value
			}
		}

		if caller.err.message != nil {
			status = caller.err.status.Default

			if caller.err.status.Reference != nil {
				ref := store.Load(caller.err.status.Reference.Resource, caller.err.status.Reference.Path)
				status = ref.Value
			}
		}

		if caller.err.codec != nil {
			err := caller.err.codec.Unmarshal(reader, store)
			if err != nil {
				return err
			}
		}

		if caller.err.functions != nil {
			err := ExecuteFunctions(caller.err.functions, store)
			if err != nil {
				return err
			}
		}

		if caller.err.metadata != nil {
			caller.response.metadata.Unmarshal(w.Header(), store)
		}
	}

	if status == nil {
		status = int64(w.Status())
	}

	if message == nil {
		message = w.Message()
	}

	store.StoreValue(template.ErrorResource, "status", status)
	store.StoreValue(template.ErrorResource, "message", message)

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
