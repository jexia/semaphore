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

func TestCaller(t *testing.T) {
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
		flow.NewNode(ctx, node, call, nil, nil),
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

	if len(dial.GetMethods()) != 1 {
		t.Errorf("unexpected methods %+v", dial.GetMethods())
	}

	rw := transport.NewResponseWriter(bytes.NewBuffer([]byte{}))
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
