package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jexia/maestro/internal/codec"
	"github.com/jexia/maestro/internal/codec/json"
	"github.com/jexia/maestro/internal/flow"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

func NewMockListener(t *testing.T, nodes flow.Nodes) (transport.Listener, int) {
	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)

	ctx := instance.NewContext()
	listener := NewListener(addr, nil)(ctx)

	json := json.NewConstructor()
	constructors := map[string]codec.Constructor{
		json.Name(): json,
	}

	endpoints := []*transport.Endpoint{
		{
			Request: transport.NewObject(NewSimpleMockSpecs(), nil),
			Flow:    flow.NewManager(ctx, "test", nodes, nil, nil, nil),
			Options: specs.Options{
				EndpointOption: "/",
				MethodOption:   http.MethodPost,
				CodecOption:    json.Name(),
			},
			Response: transport.NewObject(NewSimpleMockSpecs(), nil),
		},
	}

	listener.Handle(ctx, endpoints, constructors)
	return listener, port
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

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
	result, err := http.Post(endpoint, "application/json", strings.NewReader(`{"message":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", result.StatusCode)
	}

	if called != 1 {
		t.Errorf("unexpected called %d, expected %d", called, len(nodes))
	}
}

func TestListenerBadRequest(t *testing.T) {
	called := 0
	nodes := flow.Nodes{
		{
			Name:     "first",
			Previous: flow.Nodes{},
			Call: NewCallerFunc(func(ctx context.Context, refs refs.Store) error {
				called++
				return nil
			}),
			Next: flow.Nodes{},
		},
	}

	listener, port := NewMockListener(t, nodes)
	defer listener.Close()
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
	result, err := http.Post(endpoint, "application/json", strings.NewReader(`{"message":}`))
	if err != nil {
		t.Fatal(err)
	}

	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d, expected %d", result.StatusCode, http.StatusBadRequest)
	}

	if called == 1 {
		t.Errorf("unexpected called %d, expected %d", called, 0)
	}
}

func TestPathReferences(t *testing.T) {
	message := "active"
	nodes := flow.Nodes{
		{
			Name:     "first",
			Previous: flow.Nodes{},
			Call: NewCallerFunc(func(ctx context.Context, refs refs.Store) error {
				ref := refs.Load("input", "message")
				if ref == nil {
					t.Fatal("input:message ref has not been set")
				}

				if ref.Value != message {
					t.Fatalf("unexpected ref value %+v, expected %+v", ref.Value, message)
				}

				return nil
			}),
			Next: flow.Nodes{},
		},
	}

	listener, port := NewMockListener(t, nodes)
	defer listener.Close()

	ctx := instance.NewContext()
	endpoints := []*transport.Endpoint{
		{
			Flow: flow.NewManager(ctx, "test", nodes, nil, nil, nil),
			Options: specs.Options{
				"endpoint": "/:message",
				"method":   "GET",
			},
		},
	}

	listener.Handle(ctx, endpoints, nil)
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/"+message, port)
	http.Get(endpoint)
}

// func TestListenerForwarding(t *testing.T) {
// 	ctx := instance.NewContext()

// 	mock := fmt.Sprintf(":%d", AvailablePort(t))
// 	forward := fmt.Sprintf(":%d", AvailablePort(t))

// 	forwarded := 0

// 	go http.ListenAndServe(forward, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("forwarder")
// 		// set-up a simple forward server which always returns a 200
// 		forwarded++
// 		return
// 	}))

// 	listener := NewListener(mock, nil)(ctx)

// 	json := json.NewConstructor()
// 	constructors := map[string]codec.Constructor{
// 		json.Name(): json,
// 	}

// 	endpoints := []*transport.Endpoint{
// 		{
// 			Flow: flow.NewManager(ctx, "test", nil, nil, nil, nil),
// 			Options: specs.Options{
// 				EndpointOption: "/",
// 				MethodOption:   http.MethodPost,
// 				CodecOption:    json.Name(),
// 			},
// 			Forward: &transport.Forward{
// 				Service: &specs.Service{
// 					Host: fmt.Sprintf("http://127.0.0.1%s", forward),
// 				},
// 			},
// 		},
// 	}

// 	listener.Handle(ctx, endpoints, constructors)
// 	defer listener.Close()
// 	go listener.Serve()

// 	// Some CI pipelines take a little while before the listener is active
// 	time.Sleep(100 * time.Millisecond)

// 	endpoint := fmt.Sprintf("http://127.0.0.1%s/", mock)
// 	http.Get(endpoint)

// 	if forwarded != 1 {
// 		t.Fatalf("unexpected counter result %d, expected service request counter to be 1", forwarded)
// 	}
// }
