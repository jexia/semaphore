package flow

import (
	"context"
	"errors"
	"io"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/header"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	log "github.com/sirupsen/logrus"
)

// NewRequest constructs a new request for the given codec and header manager
func NewRequest(codec codec.Manager, header *header.Manager) *Request {
	return &Request{
		codec:  codec,
		header: header,
	}
}

// NewCall constructs a new flow caller from the given protocol caller and
func NewCall(node *specs.Node, protocol protocol.Call, request *Request, response *Request) Call {
	return &Caller{
		node:     node,
		protocol: protocol,
		request:  request,
		response: response,
	}
}

// Request represents a codec and header manager
type Request struct {
	codec  codec.Manager
	header *header.Manager
}

// Caller represents a flow protocol caller
type Caller struct {
	node     *specs.Node
	protocol protocol.Call
	request  *Request
	response *Request
}

// References returns the references inside the configured protocol caller
func (caller *Caller) References() []*specs.Property {
	return caller.protocol.References()
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
		Context: ctx,
		Body:    body,
		Header:  caller.request.header.Marshal(store),
	}

	go func() {
		defer writer.Close()
		err := caller.protocol.Call(w, r, store)
		if err != nil {
			log.Error(err)
		}
	}()

	err = caller.response.codec.Unmarshal(reader, store)
	if err != nil {
		return err
	}

	if !protocol.StatusSuccess(w.Status()) {
		log.WithFields(log.Fields{
			"node":   caller.node.GetName(),
			"status": w.Status(),
		}).Error("Faulty status code")

		return errors.New("unexpected status code, rollback required")
	}

	return nil
}
