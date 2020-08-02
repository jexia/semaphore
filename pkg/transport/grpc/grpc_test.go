package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

type caller struct {
	fn func(context.Context, references.Store) error
}

func (caller *caller) Do(ctx context.Context, store references.Store) error {
	return caller.fn(ctx, store)
}

func (caller *caller) References() []*specs.Property {
	return nil
}

func NewCallerFunc(fn func(context.Context, references.Store) error) flow.Call {
	return &caller{fn: fn}
}

func NewMockListener(t *testing.T, nodes flow.Nodes, errs transport.Errs) (transport.Listener, int) {
	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)

	ctx := instance.NewContext()
	listener := NewListener(addr, nil)(ctx)

	constructors := map[string]codec.Constructor{}
	endpoints := []*transport.Endpoint{
		{
			Request: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
			Flow:    flow.NewManager(ctx, "test", nodes, nil, nil, nil),
			Options: specs.Options{
				ServiceOption: "mock",
				MethodOption:  "simple",
				PackageOption: "pkg",
			},
			Errs:     errs,
			Response: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
		},
	}

	listener.Handle(ctx, endpoints, constructors)
	go listener.Serve()

	return listener, port
}

func NewSimpleMockSpecs() *specs.ParameterMap {
	return &specs.ParameterMap{
		Property: &specs.Property{
			Type:  types.Message,
			Label: labels.Optional,
			Nested: map[string]*specs.Property{
				"message": {
					Comment:  "mock",
					Position: 1,
					Name:     "message",
					Path:     "message",
					Type:     types.String,
					Label:    labels.Optional,
				},
			},
		},
	}
}

func AvailablePort(t *testing.T) int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
