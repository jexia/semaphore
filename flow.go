package maestro

import (
	"context"
	"errors"
	"io"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	log "github.com/sirupsen/logrus"
)

// ConstructFlowManager constructs the flow managers from the given specs manifest
func ConstructFlowManager(manifest *specs.Manifest, options Options) ([]*protocol.Endpoint, error) {
	endpoints := make([]*protocol.Endpoint, len(manifest.Endpoints))

	for index, endpoint := range manifest.Endpoints {
		current := manifest.GetFlow(endpoint.Flow)
		if current == nil {
			continue
		}

		nodes := make([]*flow.Node, len(current.GetNodes()))

		result := &protocol.Endpoint{
			Listener: endpoint.Listener,
			Options:  endpoint.Options,
		}

		for index, node := range current.GetNodes() {
			caller, err := ConstructCall(manifest, node, node.Call, options, current)
			if err != nil {
				return nil, err
			}

			rollback, err := ConstructCall(manifest, node, node.Rollback, options, current)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(node, caller, rollback)
		}

		collection, has := options.Codec[endpoint.Codec]
		if !has {
			return nil, trace.New(trace.WithMessage("unknown endpoint codec %s", endpoint.Codec))
		}

		if current.GetInput() != nil {
			req, err := collection.New(specs.InputResource, current.GetInput())
			if err != nil {
				return nil, err
			}

			result.Request = req
		}

		if current.GetOutput() != nil {
			res, err := collection.New(specs.InputResource, current.GetOutput())
			if err != nil {
				return nil, err
			}

			result.Response = res
			result.Header = protocol.NewHeaderManager(specs.InputResource, current.GetOutput())
		}

		forward, err := ConstructForward(manifest, current.GetForward(), options)
		if err != nil {
			return nil, err
		}

		result.Forward = forward
		result.Flow = flow.NewManager(current.GetName(), nodes)

		endpoints[index] = result
	}

	err := ConstructListeners(endpoints, options)
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}

// NewCall constructs a new flow caller from the given protocol caller and
func NewCall(node *specs.Node, protocol protocol.Call, header *protocol.HeaderManager, req codec.Manager, res codec.Manager) flow.Call {
	return &caller{
		node:     node,
		protocol: protocol,
		req:      req,
		res:      res,
	}
}

type caller struct {
	node     *specs.Node
	protocol protocol.Call
	req      codec.Manager
	res      codec.Manager
	header   *protocol.HeaderManager
}

func (caller *caller) References() []*specs.Property {
	return caller.protocol.References()
}

func (caller *caller) Do(ctx context.Context, store *refs.Store) error {
	body, err := caller.req.Marshal(store)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	w := protocol.NewResponseWriter(writer)
	r := &protocol.Request{
		Context: ctx,
		Body:    body,
		Header:  caller.header.Marshal(store),
	}

	go func() {
		defer writer.Close()
		err := caller.protocol.Call(w, r, store)
		if err != nil {
			log.Println(err)
		}
	}()

	err = caller.res.Unmarshal(reader, store)
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

// ConstructCall constructs a flow caller for the given node call.
func ConstructCall(manifest *specs.Manifest, node *specs.Node, call *specs.Call, options Options, flow specs.FlowManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	constructor := options.Callers.Get(service.GetProtocol())
	codec := options.Codec[service.GetCodec()]

	req, err := codec.New(node.GetName(), call.GetRequest())
	if err != nil {
		return nil, err
	}

	res, err := codec.New(node.GetName(), call.GetResponse())
	if err != nil {
		return nil, err
	}

	header := protocol.NewHeaderManager(node.GetName(), call.GetRequest())
	schema := options.Schema.GetService(service.GetName())

	protocol, err := constructor.New(schema, call.GetMethod(), options.Functions, service.GetOptions())
	if err != nil {
		return nil, err
	}

	caller := NewCall(node, protocol, header, req, res)
	err = strict.DefineCaller(node, manifest, protocol, flow)
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// ConstructForward constructs a flow caller for the given call.
func ConstructForward(manifest *specs.Manifest, call *specs.Call, options Options) (protocol.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.GetService())
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	schema := options.Schema.GetService(service.GetName())
	constructor := options.Callers.Get(service.GetProtocol())
	caller, err := constructor.New(schema, "", options.Functions, service.GetOptions())
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// ConstructListeners constructs the listeners from the given collection of endpoints
func ConstructListeners(endpoints []*protocol.Endpoint, options Options) error {
	collections := make(map[string][]*protocol.Endpoint, len(options.Listeners))

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		listener := options.Listeners.Get(endpoint.Listener)
		if listener == nil {
			return trace.New(trace.WithMessage("unknown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		listener := options.Listeners.Get(key)
		err := listener.Handle(collection)
		if err != nil {
			return err
		}
	}

	return nil
}
