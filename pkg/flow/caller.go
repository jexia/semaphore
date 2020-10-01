package flow

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// ErrAbortFlow represents the error thrown when a flow has to be aborted
var ErrAbortFlow = errors.New("abort flow")

// Request represents a codec and metadata manager
type Request struct {
	Functions functions.Stack
	Codec     codec.Manager
	Metadata  *metadata.Manager
}

// CallOptions represents the available options that could be used to construct a new flow caller
type CallOptions struct {
	Transport      transport.Call
	Method         transport.Method
	Request        *Request
	Response       *Request
	Err            *OnError
	ExpectedStatus []int
}

// Call represents a transport caller implementation
type Call interface {
	Do(context.Context, references.Store) error
}

// NewCall constructs a new flow caller from the given transport caller and options
func NewCall(parent *broker.Context, node *specs.Node, options *CallOptions) Call {
	if node == nil || options == nil {
		return nil
	}

	module := broker.WithModule(parent, "caller", node.ID)
	ctx := logger.WithFields(logger.WithLogger(module), zap.String("node", node.ID))

	result := &caller{
		ctx:            ctx,
		node:           node,
		transport:      options.Transport,
		method:         options.Method,
		request:        options.Request,
		response:       options.Response,
		ExpectedStatus: make(map[int]struct{}, len(options.ExpectedStatus)),
		err:            options.Err,
	}

	for _, status := range options.ExpectedStatus {
		result.ExpectedStatus[status] = struct{}{}
	}

	if len(result.ExpectedStatus) == 0 {
		result.ExpectedStatus = map[int]struct{}{transport.StatusOK: {}}
	}

	return result
}

// Caller represents a flow transport caller
type caller struct {
	ctx            *broker.Context
	node           *specs.Node
	method         transport.Method
	transport      transport.Call
	request        *Request
	response       *Request
	err            *OnError
	ExpectedStatus map[int]struct{}
}

// Do is called by the flow manager to call the configured service
func (caller *caller) Do(ctx context.Context, store references.Store) error {
	reader, writer := io.Pipe()
	w := transport.NewResponseWriter(writer)
	r := &transport.Request{
		Method: caller.method,
	}

	if caller.request != nil {
		if caller.request.Functions != nil {
			err := ExecuteFunctions(caller.request.Functions, store)
			if err != nil {
				return err
			}
		}

		if caller.request.Metadata != nil {
			r.Header = caller.request.Metadata.Marshal(store)
		}

		if caller.request.Codec != nil {
			body, err := caller.request.Codec.Marshal(store)
			if err != nil {
				return err
			}

			r.RequestCodec = caller.request.Codec.Name()
			r.Body = body
		}
	}

	if caller.response != nil {
		if caller.response.Codec != nil {
			r.ResponseCodec = caller.response.Codec.Name()
		}
	}

	if caller.transport != nil {
		// SendMsg should not be called inside a separate go routine.
		// A separate go routine could be created inside the transporter to stream the returned message to the io reader
		err := caller.transport.SendMsg(ctx, w, r, store)
		if err != nil {
			logger.Error(caller.ctx, "transporter returned a unexpected error", zap.String("node", caller.node.ID), zap.Error(err))
			return err
		}
	} else {
		writer.Close()
	}

	_, expected := caller.ExpectedStatus[w.Status()]
	if caller.transport != nil && !expected {
		logger.Error(caller.ctx, "service returned a unexpected status, aborting flow", zap.Int("status", w.Status()))

		err := caller.HandleErr(w, reader, store)
		if err != nil {
			return err
		}

		return ErrAbortFlow
	}

	if caller.response != nil {
		if caller.response.Codec != nil {
			err := caller.response.Codec.Unmarshal(reader, store)
			if err != nil {
				return fmt.Errorf("failed to unmarshal response into the store: %w", err)
			}
		}

		if caller.response.Functions != nil {
			err := ExecuteFunctions(caller.response.Functions, store)
			if err != nil {
				return err
			}
		}

		if caller.response.Metadata != nil {
			caller.response.Metadata.Unmarshal(w.Header(), store)
		}
	}

	return nil
}

// HandleErr handles a thrown service error. If a error response is defined is it decoded
func (caller *caller) HandleErr(w *transport.Writer, reader io.Reader, store references.Store) error {
	var status interface{}
	var message interface{}

	if caller.err != nil {
		if caller.err.message != nil {
			if caller.err.message.Scalar != nil {
				message = caller.err.message.Scalar.Default
			}

			if caller.err.message.Reference != nil {
				ref := store.Load(caller.err.message.Reference.Resource, caller.err.message.Reference.Path)
				if ref != nil {
					message = ref.Value
				}
			}
		}

		if caller.err.status != nil {
			if caller.err.status.Scalar != nil {
				status = caller.err.status.Scalar.Default
			}

			if caller.err.status.Reference != nil {
				ref := store.Load(caller.err.status.Reference.Resource, caller.err.status.Reference.Path)
				if ref != nil {
					status = ref.Value
				}
			}
		}

		if caller.err.codec != nil {
			err := caller.err.codec.Unmarshal(reader, store)
			if err != nil {
				return err
			}
		}

		if caller.err.stack != nil {
			err := ExecuteFunctions(caller.err.stack, store)
			if err != nil {
				return err
			}
		}

		if caller.err.metadata != nil {
			caller.response.Metadata.Unmarshal(w.Header(), store)
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
func ExecuteFunctions(stack functions.Stack, store references.Store) error {
	for key, function := range stack {
		resource := template.JoinPath(template.StackResource, key)
		err := function.Fn(references.NewPrefixStore(store, resource, ""))
		if err != nil {
			return err
		}
	}

	return nil
}
