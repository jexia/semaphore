package http

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

func TestListener(t *testing.T) {
	called := 0
	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)
	listener := NewListener(addr, nil)

	defer listener.Close()

	nodes := flow.Nodes{
		{
			Name:     "first",
			Previous: flow.Nodes{},
			Call: func(ctx context.Context, refs *refs.Store) error {
				called++
				return nil
			},
			Next: flow.Nodes{},
		},
	}

	endpoints := []*protocol.Endpoint{
		{
			Flow: flow.NewManager("test", nodes),
			Options: specs.Options{
				"endpoint": "/",
				"method":   "GET",
			},
		},
	}

	listener.Handle(endpoints)
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
	http.Get(endpoint)

	if called != 1 {
		t.Errorf("unexpected called %d, expected %d", called, len(nodes))
	}
}

func TestPathReferences(t *testing.T) {
	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)
	listener := NewListener(addr, nil)

	defer listener.Close()

	message := "active"
	nodes := flow.Nodes{
		{
			Name:     "first",
			Previous: flow.Nodes{},
			Call: func(ctx context.Context, refs *refs.Store) error {
				ref := refs.Load("input", "message")
				if ref == nil {
					t.Fatal("input:message ref has not been set")
				}

				if ref.Value != message {
					t.Fatalf("unexpected ref value %+v, expected %+v", ref.Value, message)
				}

				return nil
			},
			Next: flow.Nodes{},
		},
	}

	endpoints := []*protocol.Endpoint{
		{
			Flow: flow.NewManager("test", nodes),
			Options: specs.Options{
				"endpoint": "/:message",
				"method":   "GET",
			},
		},
	}

	listener.Handle(endpoints)
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/"+message, port)
	http.Get(endpoint)
}
