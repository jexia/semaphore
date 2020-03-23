package flow

import (
	"context"
	"io"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/sirupsen/logrus"
)

// NewRequest constructs a new request for the given codec and header manager
func NewRequest(codec codec.Manager, metadata *metadata.Manager) *Request {
	return &Request{
		codec:    codec,
		metadata: metadata,
	}
}

// NewCall constructs a new flow caller from the given protocol caller and
func NewCall(ctx context.Context, node *specs.Node, protocol protocol.Call, method string, request *Request, response *Request) Call {
	return &Caller{
		ctx:      ctx,
		node:     node,
		protocol: protocol,
		method:   protocol.GetMethod(method),
		request:  request,
		response: response,
	}
}

// Request represents a codec and header manager
type Request struct {
	codec    codec.Manager
	metadata *metadata.Manager
}

// Caller represents a flow protocol caller
type Caller struct {
	ctx      context.Context
	node     *specs.Node
	method   protocol.Method
	protocol protocol.Call
	request  *Request
	response *Request
}

// References returns the references inside the configured protocol caller
func (caller *Caller) References() []*specs.Property {
	if caller.method == nil {
		return make([]*specs.Property, 0)
	}

	return caller.method.References()
}

// Do is called by the flow manager to call the configured service
func (caller *Caller) Do(ctx context.Context, store *refs.Store) error {
	body, err := caller.request.codec.Marshal(store)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	w := protocol.NewResponseWriter(writer)
	r := &protocol.Request{
		Header: caller.request.metadata.Marshal(store),
		Method: caller.method,
		Body:   body,
	}

	result := make(chan error, 1)
	defer close(result)

	go func() {
		defer writer.Close()
		result <- caller.protocol.SendMsg(ctx, w, r, store)
	}()

	err = caller.response.codec.Unmarshal(reader, store)
	if err != nil {
		return err
	}

	err = <-result
	if err != nil {
		logger.FromCtx(caller.ctx, logger.Flow).WithFields(logrus.Fields{
			"node": caller.node.GetName(),
			"err":  err,
		}).Error("Service error")

		return err
	}

	caller.response.metadata.Unmarshal(w.Header(), store)

	return nil
}
