package flow

import (
	"context"
	"io"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
	"github.com/sirupsen/logrus"
)

// NewRequest constructs a new request for the given codec and header manager
func NewRequest(codec codec.Manager, metadata *metadata.Manager) *Request {
	return &Request{
		codec:    codec,
		metadata: metadata,
	}
}

// NewCall constructs a new flow caller from the given transport caller and
func NewCall(ctx instance.Context, node *specs.Node, transport transport.Call, method string, request *Request, response *Request) Call {
	return &Caller{
		ctx:       ctx,
		node:      node,
		transport: transport,
		method:    transport.GetMethod(method),
		request:   request,
		response:  response,
	}
}

// Request represents a codec and header manager
type Request struct {
	codec    codec.Manager
	metadata *metadata.Manager
}

// Caller represents a flow transport caller
type Caller struct {
	ctx       instance.Context
	node      *specs.Node
	method    transport.Method
	transport transport.Call
	request   *Request
	response  *Request
}

// References returns the references inside the configured transport caller
func (caller *Caller) References() []*specs.Property {
	if caller.method == nil {
		return make([]*specs.Property, 0)
	}

	return caller.method.References()
}

// Do is called by the flow manager to call the configured service
func (caller *Caller) Do(ctx context.Context, store *specs.Store) error {
	body, err := caller.request.codec.Marshal(store)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	w := transport.NewResponseWriter(writer)
	r := &transport.Request{
		Header: caller.request.metadata.Marshal(store),
		Method: caller.method,
		Body:   body,
	}

	result := make(chan error, 1)
	defer close(result)

	go func() {
		defer writer.Close()
		result <- caller.transport.SendMsg(ctx, w, r, store)
	}()

	err = caller.response.codec.Unmarshal(reader, store)
	if err != nil {
		return err
	}

	err = <-result
	if err != nil {
		caller.ctx.Logger(logger.Flow).WithFields(logrus.Fields{
			"node": caller.node.GetName(),
			"err":  err,
		}).Error("Service error")

		return err
	}

	caller.response.metadata.Unmarshal(w.Header(), store)

	return nil
}
