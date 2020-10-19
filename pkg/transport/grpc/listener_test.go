package grpc

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

func TestNewListener(t *testing.T) {
	tests := map[string]func(*testing.T){
		"simple": func(t *testing.T) {
			constructor := NewListener(":0", specs.Options{})
			if constructor == nil {
				t.Fatal("nil listener constructor")
			}

			ctx := logger.WithLogger(broker.NewBackground())
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

func TestErrorHandlingListener(t *testing.T) {
	type test struct {
		caller   func(references.Store)
		err      *specs.OnError
		expected int
		result   string
	}

	tests := map[string]test{
		"simple": {
			err: &specs.OnError{
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.Int64,
							Default: int64(500),
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.String,
							Default: "database broken",
						},
					},
				},
			},
			expected: 500,
			result:   "database broken",
		},
		"reference": {
			caller: func(store references.Store) {
				store.StoreValue("error", "status", int64(429))
				store.StoreValue("error", "message", "reference value")
			},
			err: &specs.OnError{
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.Int64,
						},
						Reference: &specs.PropertyReference{
							Resource: "error",
							Path:     "status",
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
						Reference: &specs.PropertyReference{
							Resource: "error",
							Path:     "message",
						},
					},
				},
			},
			expected: 429,
			result:   "reference value",
		},
		"input": {
			caller: func(store references.Store) {
				store.StoreValue("error", "status", int64(429))
				store.StoreValue("input", "message", "reference value")
			},
			err: &specs.OnError{
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.Int64,
						},
						Reference: &specs.PropertyReference{
							Resource: "error",
							Path:     "status",
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "message",
						},
					},
				},
			},
			expected: 429,
			result:   "reference value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			node := &specs.Node{
				ID:      "first",
				OnError: test.err,
			}

			called := 0
			call := NewCallerFunc(func(ctx context.Context, refs references.Store) error {
				called++

				if test.caller != nil {
					test.caller(refs)
				}

				return flow.ErrAbortFlow
			})

			nodes := flow.Nodes{
				flow.NewNode(ctx, node, flow.WithCall(call)),
			}

			obj := transport.NewObject(node.OnError.Response, node.OnError.Status, node.OnError.Message)
			errs := transport.Errs{
				node.OnError: obj,
			}

			listener, port := NewMockListener(t, nodes, errs)
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

			rw := transport.NewResponseWriter(&DiscardWriter{})
			rq := &transport.Request{
				Header: metadata.MD{},
				Method: dial.GetMethod("simple"),
				Body:   bytes.NewBuffer([]byte{}),
			}

			err = dial.SendMsg(context.Background(), rw, rq, references.NewReferenceStore(0))
			if err != nil {
				t.Fatalf("unrecoverable err returned '%s'", err)
			}

			if called != 1 {
				t.Errorf("unexpected called %d, expected %d", called, len(nodes))
			}

			if rw.Status() != test.expected {
				t.Fatalf("unexpected status %d, expected %d", rw.Status(), test.expected)
			}

			if rw.Message() != test.result {
				t.Fatalf("unexpected message %s, expected %s", rw.Message(), test.result)
			}
		})
	}
}
