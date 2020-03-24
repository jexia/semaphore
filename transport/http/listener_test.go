package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
)

func NewMockListener(t *testing.T, nodes flow.Nodes) (transport.Listener, int) {
	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)
	listener := NewListener(addr, nil)

	ctx := context.Background()
	ctx = logger.WithValue(ctx)
	listener.Context(ctx)

	json := json.NewConstructor()
	constructors := map[string]codec.Constructor{
		json.Name(): json,
	}

	endpoints := []*transport.Endpoint{
		{
			Request: NewSimpleMockSpecs(),
			Flow:    flow.NewManager(ctx, "test", nodes),
			Options: specs.Options{
				EndpointOption: "/",
				MethodOption:   http.MethodPost,
				CodecOption:    json.Name(),
			},
			Response: NewSimpleMockSpecs(),
		},
	}

	listener.Handle(endpoints, constructors)
	return listener, port
}

func TestListener(t *testing.T) {
	ctx := context.Background()
	ctx = logger.WithValue(ctx)

	specs := &specs.Node{
		Name: "first",
	}

	called := 0
	call := NewCallerFunc(func(ctx context.Context, refs *refs.Store) error {
		called++
		return nil
	})

	nodes := flow.Nodes{
		flow.NewNode(ctx, specs, call, nil),
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
			Call: NewCallerFunc(func(ctx context.Context, refs *refs.Store) error {
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
			Call: NewCallerFunc(func(ctx context.Context, refs *refs.Store) error {
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

	ctx := context.Background()
	ctx = logger.WithValue(ctx)

	endpoints := []*transport.Endpoint{
		{
			Flow: flow.NewManager(ctx, "test", nodes),
			Options: specs.Options{
				"endpoint": "/:message",
				"method":   "GET",
			},
		},
	}

	listener.Handle(endpoints, nil)
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/"+message, port)
	http.Get(endpoint)
}
