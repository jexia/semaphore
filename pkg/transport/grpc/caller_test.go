package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

type DiscardWriter struct {
}

func (d *DiscardWriter) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func (d *DiscardWriter) Close() error {
	return nil
}

func TestUnknownMethod(t *testing.T) {
	type fields struct {
		Method string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Method: "get"},
			"unknown service method get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUnknownMethod{
				Method: "get",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCaller(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	node := &specs.Node{
		ID: "first",
	}

	called := 0
	call := NewCallerFunc(func(ctx context.Context, refs references.Store) error {
		called++
		return nil
	})

	nodes := flow.Nodes{
		flow.NewNode(ctx, node, flow.WithCall(call)),
	}

	listener, port := NewMockListener(t, nodes, nil)
	defer listener.Close()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	constructor := NewCaller()
	caller := constructor(ctx)

	service := &specs.Service{
		Name:          "mock",
		Package:       "pkg",
		Host:          fmt.Sprintf("127.0.0.1:%d", port),
		Transport:     "grpc",
		RequestCodec:  "proto",
		ResponseCodec: "proto",
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

	rw := transport.NewResponseWriter(&DiscardWriter{})
	rq := &transport.Request{
		Header: metadata.MD{},
		Method: dial.GetMethod("simple"),
		Body:   bytes.NewBuffer([]byte{}),
	}

	err = dial.SendMsg(context.Background(), rw, rq, references.NewReferenceStore(0))
	if err != nil {
		t.Fatal(err)
	}

	if called != 1 {
		t.Errorf("unexpected called %d, expected %d", called, len(nodes))
	}
}
