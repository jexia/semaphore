package grpc

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jexia/maestro/pkg/flow"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

func TestNewListener(t *testing.T) {
	tests := map[string]func(*testing.T){
		"simple": func(t *testing.T) {
			constructor := NewListener(":0", specs.Options{})
			if constructor == nil {
				t.Fatal("nil listener constructor")
			}

			ctx := instance.NewContext()
			listener := constructor(ctx)
			if constructor == nil {
				t.Fatal("nil listener")
			}

			if listener.Name() != "grpc" {
				t.Fatal("unknown listener name")
			}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test(t)
		})
	}
}

func TestListener(t *testing.T) {
	ctx := instance.NewContext()
	node := &specs.Node{
		Name: "first",
	}

	called := 0
	call := NewCallerFunc(func(ctx context.Context, refs refs.Store) error {
		called++
		return nil
	})

	nodes := flow.Nodes{
		flow.NewNode(ctx, node, nil, call, nil, nil),
	}

	listener, port := NewMockListener(t, nodes)

	defer listener.Close()
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	constructor := NewCaller()
	caller := constructor(ctx)

	service := &specs.Service{
		Name:      "mock",
		Package:   "pkg",
		Host:      fmt.Sprintf("127.0.0.1:%d", port),
		Transport: "grpc",
		Codec:     "proto",
		Methods: []*specs.Method{
			{
				Name:    "simple",
				Options: specs.Options{},
			},
		},
		Options: specs.Options{},
	}

	dial, err := caller.Dial(service, nil, specs.Options{})
	if err != nil {
		t.Fatal(err)
	}

	defer dial.Close()

	rw := transport.NewResponseWriter(&DiscardWriter{})
	rq := &transport.Request{
		Header: metadata.MD{},
		Method: dial.GetMethod("simple"),
		Body:   bytes.NewBuffer([]byte{}),
	}

	err = dial.SendMsg(context.Background(), rw, rq, refs.NewReferenceStore(0))
	if err != nil {
		t.Fatal(err)
	}

	if called != 1 {
		t.Errorf("unexpected called %d, expected %d", called, len(nodes))
	}
}
