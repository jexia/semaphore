package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/flow"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jexia/maestro/pkg/transport"
)

type caller struct {
	fn func(context.Context, refs.Store) error
}

func (caller *caller) Do(ctx context.Context, store refs.Store) error {
	return caller.fn(ctx, store)
}

func (caller *caller) References() []*specs.Property {
	return nil
}

func NewCallerFunc(fn func(context.Context, refs.Store) error) flow.Call {
	return &caller{fn: fn}
}

func NewMockListener(t *testing.T, nodes flow.Nodes) (transport.Listener, int) {
	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)

	ctx := instance.NewContext()
	listener := NewListener(addr, nil)(ctx)

	constructors := map[string]codec.Constructor{}
	endpoints := []*transport.Endpoint{
		{
			Request: NewSimpleMockSpecs(),
			Flow:    flow.NewManager(ctx, "test", nodes, nil),
			Options: specs.Options{
				ServiceOption: "mock",
				MethodOption:  "simple",
				PackageOption: "pkg",
			},
			Response: NewSimpleMockSpecs(),
		},
	}

	listener.Handle(endpoints, constructors)
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
