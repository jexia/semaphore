package http

import (
	"bytes"
	"context"
	ejson "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/refs"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

func NewMockListener(t *testing.T, nodes flow.Nodes, errs transport.Errs) (transport.Listener, string) {
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
			Request: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
			Flow:    flow.NewManager(ctx, "test", nodes, nil, nil, nil),
			Errs:    errs,
			Options: specs.Options{
				EndpointOption: "/",
				MethodOption:   http.MethodPost,
				CodecOption:    json.Name(),
			},
			Response: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
		},
	}

	listener.Handle(ctx, endpoints, constructors)
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
	return listener, endpoint
}

func TestListener(t *testing.T) {
	ctx := instance.NewContext()
	node := &specs.Node{
		ID: "first",
	}

	called := 0
	call := NewCallerFunc(func(ctx context.Context, refs refs.Store) error {
		called++
		return nil
	})

	nodes := flow.Nodes{
		flow.NewNode(ctx, node, nil, call, nil, nil),
	}

	listener, endpoint := NewMockListener(t, nodes, nil)
	defer listener.Close()

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

	listener, endpoint := NewMockListener(t, nodes, nil)
	defer listener.Close()

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

	listener, port := NewMockListener(t, nodes, nil)
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

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/"+message, port)
	http.Get(endpoint)
}

func TestStoringParams(t *testing.T) {
	ctx := instance.NewContext()
	node := &specs.Node{
		ID: "first",
	}

	path := "message"
	expected := "sample"
	called := 0

	call := NewCallerFunc(func(ctx context.Context, refs refs.Store) error {
		ref := refs.Load("input", path)
		if ref == nil {
			t.Fatal("reference not set")
		}

		if ref.Value == nil {
			t.Fatal("reference value not set")
		}

		if ref.Value != expected {
			t.Fatalf("unexpected value '%+v', expected '%s'", ref.Value, expected)
		}

		called++
		return nil
	})

	nodes := flow.Nodes{
		flow.NewNode(ctx, node, nil, call, nil, nil),
	}

	listener, endpoint := NewMockListener(t, nodes, nil)
	defer listener.Close()

	uri, err := url.Parse(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	query := uri.Query()
	query.Add(path, expected)
	uri.RawQuery = query.Encode()

	t.Log(uri.String())

	res, err := http.Post(uri.String(), "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d, expected %d", res.StatusCode, http.StatusOK)
	}

	if called != 1 {
		t.Fatalf("unexpected counter result %d, expected service request counter to be 1", called)
	}
}

func TestListenerForwarding(t *testing.T) {
	ctx := instance.NewContext()

	mock := fmt.Sprintf(":%d", AvailablePort(t))
	forward := fmt.Sprintf(":%d", AvailablePort(t))

	forwarded := 0

	go http.ListenAndServe(forward, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set-up a simple forward server which always returns a 200
		forwarded++
		return
	}))

	listener := NewListener(mock, nil)(ctx)

	json := json.NewConstructor()
	constructors := map[string]codec.Constructor{
		json.Name(): json,
	}

	endpoints := []*transport.Endpoint{
		{
			Flow: flow.NewManager(ctx, "test", nil, nil, nil, nil),
			Options: specs.Options{
				EndpointOption: "/",
				MethodOption:   http.MethodGet,
				CodecOption:    json.Name(),
			},
			Forward: &transport.Forward{
				Service: &specs.Service{
					Host: fmt.Sprintf("http://127.0.0.1%s", forward),
				},
			},
		},
	}

	listener.Handle(ctx, endpoints, constructors)
	defer listener.Close()
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1%s/", mock)
	res, err := http.Get(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d, expected %d", res.StatusCode, http.StatusOK)
	}

	if forwarded != 1 {
		t.Fatalf("unexpected counter result %d, expected service request counter to be 1", forwarded)
	}
}

func TestListenerErrorHandling(t *testing.T) {
	type test struct {
		input    map[string]string
		caller   func(refs.Store)
		err      *specs.OnError
		expected int
		response map[string]interface{}
	}

	tests := map[string]test{
		"simple": {
			input: map[string]string{
				"message": "value",
			},
			caller: func(store refs.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(500))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Type:  types.Message,
						Label: labels.Optional,
						Nested: map[string]*specs.Property{
							"status": {
								Name:  "status",
								Path:  "status",
								Type:  types.Int64,
								Label: labels.Optional,
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "status",
								},
							},
							"message": {
								Name:  "message",
								Path:  "message",
								Type:  types.String,
								Label: labels.Optional,
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "message",
								},
							},
						},
					},
				},
				Status: &specs.Property{
					Type:    types.Int64,
					Label:   labels.Optional,
					Default: int64(500),
				},
				Message: &specs.Property{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "value",
				},
			},
			expected: 500,
			response: map[string]interface{}{
				"status":  500,
				"message": "value",
			},
		},
		"reference": {
			input: map[string]string{
				"message": "value",
			},
			caller: func(store refs.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(500))
				store.StoreValue("input", "status", int64(401))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Type:  types.Message,
						Label: labels.Optional,
						Nested: map[string]*specs.Property{
							"status": {
								Name:  "status",
								Path:  "status",
								Type:  types.Int64,
								Label: labels.Optional,
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "status",
								},
							},
							"message": {
								Name:  "message",
								Path:  "message",
								Type:  types.String,
								Label: labels.Optional,
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "message",
								},
							},
						},
					},
				},
				Status: &specs.Property{
					Type:    types.Int64,
					Label:   labels.Optional,
					Default: int64(401),
				},
				Message: &specs.Property{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "value",
				},
			},

			expected: 401,
			response: map[string]interface{}{
				"status":  401,
				"message": "value",
			},
		},
		"not_found": {
			input: map[string]string{
				"message": "value",
			},
			caller: func(store refs.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(404))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Type:  types.Message,
						Label: labels.Optional,
						Nested: map[string]*specs.Property{
							"status": {
								Name:  "status",
								Path:  "status",
								Type:  types.Int64,
								Label: labels.Optional,
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "status",
								},
							},
							"message": {
								Name:  "message",
								Path:  "message",
								Type:  types.String,
								Label: labels.Optional,
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "message",
								},
							},
						},
					},
				},
				Status: &specs.Property{
					Type:    types.Int64,
					Label:   labels.Optional,
					Default: int64(404),
				},
				Message: &specs.Property{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "value",
				},
			},

			expected: 404,
			response: map[string]interface{}{
				"status":  404,
				"message": "value",
			},
		},
		"complex": {
			input: map[string]string{
				"message": "value",
			},
			caller: func(store refs.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(404))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Type:  types.Message,
						Label: labels.Optional,
						Nested: map[string]*specs.Property{
							"meta": {
								Name:  "meta",
								Path:  "meta",
								Type:  types.Message,
								Label: labels.Optional,
								Nested: map[string]*specs.Property{
									"status": {
										Name:  "status",
										Path:  "meta.status",
										Type:  types.Int64,
										Label: labels.Optional,
										Reference: &specs.PropertyReference{
											Resource: "error",
											Path:     "status",
										},
									},
									"message": {
										Name:  "message",
										Path:  "meta.message",
										Type:  types.String,
										Label: labels.Optional,
										Reference: &specs.PropertyReference{
											Resource: "error",
											Path:     "message",
										},
									},
								},
							},
							"const": {
								Name:    "const",
								Path:    "const",
								Type:    types.String,
								Label:   labels.Optional,
								Default: "custom message",
							},
						},
					},
				},
				Status: &specs.Property{
					Type:    types.Int64,
					Label:   labels.Optional,
					Default: int64(404),
				},
				Message: &specs.Property{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "value",
				},
			},

			expected: 404,
			response: map[string]interface{}{
				"meta": map[string]interface{}{
					"status":  404,
					"message": "value",
				},
				"const": "custom message",
			},
		},
		"empty": {
			input: map[string]string{
				"message": "value",
			},
			caller: func(store refs.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(404))
			},
			err: &specs.OnError{
				Response: nil,
				Status: &specs.Property{
					Type:    types.Int64,
					Label:   labels.Optional,
					Default: int64(404),
				},
				Message: &specs.Property{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "value",
				},
			},

			expected: 404,
			response: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			node := &specs.Node{
				ID:      "first",
				OnError: test.err,
			}

			called := 0
			call := NewCallerFunc(func(ctx context.Context, refs refs.Store) error {
				called++

				if test.caller != nil {
					test.caller(refs)
				}

				return flow.ErrAbortFlow
			})

			nodes := flow.Nodes{
				flow.NewNode(ctx, node, nil, call, nil, nil),
			}

			obj := transport.NewObject(node.OnError.Response, node.OnError.Status, nil)
			errs := transport.Errs{
				node.OnError: obj,
			}

			listener, endpoint := NewMockListener(t, nodes, errs)
			defer listener.Close()

			uri, err := url.Parse(endpoint)
			if err != nil {
				t.Fatal(err)
			}

			query := uri.Query()

			for key, val := range test.input {
				query.Add(key, val)
			}

			uri.RawQuery = query.Encode()

			t.Log(uri.String())

			res, err := http.Post(uri.String(), "", nil)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != test.expected {
				t.Fatalf("unexpected status code %d, expected %d", res.StatusCode, test.expected)
			}

			if called != 1 {
				t.Fatalf("unexpected counter result %d, expected service request counter to be 1", called)
			}

			expected, err := ejson.Marshal(test.response)
			if err != nil {
				t.Fatal(err)
			}

			equal, left, right, err := JSONEqual(res.Body, bytes.NewBuffer(expected))
			if err == io.EOF && test.response == nil {
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if !equal {
				t.Fatalf("unexpected response %+v, expected %+v", left, right)
			}
		})
	}
}
