package http

import (
	"bytes"
	"context"
	ejson "encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

func NewMockListener(t *testing.T, nodes flow.Nodes, errs transport.Errs) (transport.Listener, string) {
	var (
		port   = AvailablePort(t)
		addr   = fmt.Sprintf(":%d", port)
		origin = []string{"test.com"}

		ctx      = logger.WithLogger(broker.NewBackground())
		listener = NewListener(addr, WithOrigins(origin))(ctx)

		json = json.NewConstructor()

		constructors = map[string]codec.Constructor{
			json.Name(): json,
		}

		endpoints = []*transport.Endpoint{
			{
				Request: &transport.Object{
					Definition: NewSimpleMockSpecs(),
				},
				Flow: flow.NewManager(ctx, "test", nodes, nil, nil, nil),
				Errs: errs,
				Options: specs.Options{
					EndpointOption: "/",
					MethodOption:   http.MethodPost,
					CodecOption:    json.Name(),
				},
				Response: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
			},
		}
	)

	listener.Handle(ctx, endpoints, constructors)
	go listener.Serve()

	// Some CI pipelines take a little while before the listener is active
	time.Sleep(100 * time.Millisecond)

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
	return listener, endpoint
}

func TestListenerRouteConflict(t *testing.T) {
	t.Run("detect route confilcts and catch panics thrown by julienschmidt httprouter", func(t *testing.T) {
		var (
			port   = AvailablePort(t)
			addr   = fmt.Sprintf(":%d", port)
			origin = []string{"test.com"}

			ctx      = logger.WithLogger(broker.NewBackground())
			listener = NewListener(addr, WithOrigins(origin))(ctx)

			json = json.NewConstructor()

			constructors = map[string]codec.Constructor{
				json.Name(): json,
			}

			endpoints = []*transport.Endpoint{
				{
					Flow: flow.NewManager(ctx, "test", flow.Nodes{}, nil, nil, nil),
					Options: specs.Options{
						EndpointOption: "/static",
						MethodOption:   http.MethodGet,
						CodecOption:    json.Name(),
					},
				},
				{
					Flow: flow.NewManager(ctx, "test", flow.Nodes{}, nil, nil, nil),
					Options: specs.Options{
						EndpointOption: "/:variable",
						MethodOption:   http.MethodGet,
						CodecOption:    json.Name(),
					},
				},
			}
		)

		err := listener.Handle(ctx, endpoints, constructors)
		if err == nil {
			t.Errorf("error was expected")
		}

		if !errors.As(err, new(ErrRouteConflict)) {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func TestErrUndefinedCodec(t *testing.T) {
	type fields struct {
		Codec string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Codec: "json"},
			"request codec not found 'json'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedCodec{
				Codec: "json",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrInvalidHost(t *testing.T) {
	type fields struct {
		Host string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Host: "127.0.0.1"},
			"unable to parse the proxy forward host '127.0.0.1'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrInvalidHost{
				Host: "127.0.0.1",
			}
			if got := e.Prettify(); got.Message != tt.want && got.Unwrap().Error() != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	type test struct {
		headers          map[string]string
		shouldPresent    map[string]string
		shouldNotPresent []string
	}

	tests := map[string]test{
		"bad origin": {
			headers: map[string]string{"Origin": "un.known"},
			shouldPresent: map[string]string{
				"Allow": "OPTIONS, POST",
			},
			shouldNotPresent: []string{
				"Access-Control-Allow-Origin",
				"Access-Control-Allow-Headers",
			},
		},
		"good origin": {
			headers: map[string]string{
				"Origin":                         "test.com",
				"Access-Control-Request-Method":  "POST",
				"Access-Control-Request-Headers": "Authorization",
			},
			shouldPresent: map[string]string{
				"Access-Control-Allow-Origin":  "test.com",
				"Access-Control-Allow-Methods": "POST",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				ctx  = logger.WithLogger(broker.NewBackground())
				node = &specs.Node{
					ID: "first",
				}

				called = 0
				call   = NewCallerFunc(func(ctx context.Context, refs references.Store) error {
					called++
					return nil
				})

				nodes = flow.Nodes{
					flow.NewNode(ctx, node, flow.WithCall(call)),
				}
			)

			listener, endpoint := NewMockListener(t, nodes, nil)
			defer listener.Close()

			req, err := http.NewRequest(http.MethodOptions, endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			for header, value := range test.headers {
				req.Header.Set(header, value)
			}

			result, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if result.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code %d", result.StatusCode)
			}

			for header, expected := range test.shouldPresent {
				if value := result.Header.Get(header); value != expected {
					t.Errorf("header %q should have value %q instead of %q", header, expected, value)
				}
			}

			for _, header := range test.shouldNotPresent {
				if result.Header.Get(header) != "" {
					t.Errorf("header %q should not present in the response", header)
				}
			}

			if called != 0 {
				t.Errorf("node handler should not be called")
			}
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
			Call: NewCallerFunc(func(ctx context.Context, refs references.Store) error {
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
			Call: NewCallerFunc(func(ctx context.Context, refs references.Store) error {
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

	listener, host := NewMockListener(t, nodes, nil)
	defer listener.Close()

	ctx := logger.WithLogger(broker.NewBackground())
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

	endpoint := host + message
	_, err := http.Get(endpoint)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStoringParams(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	node := &specs.Node{
		ID: "first",
	}

	path := "message"
	expected := "sample"
	called := 0

	call := NewCallerFunc(func(ctx context.Context, refs references.Store) error {
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
		flow.NewNode(ctx, node, flow.WithCall(call)),
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
	ctx := logger.WithLogger(broker.NewBackground())

	mock := fmt.Sprintf(":%d", AvailablePort(t))
	forward := fmt.Sprintf(":%d", AvailablePort(t))

	forwarded := 0

	go http.ListenAndServe(forward, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set-up a simple forward server which always returns a 200
		forwarded++
	}))

	listener := NewListener(mock)(ctx)

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
		caller   func(references.Store)
		err      *specs.OnError
		expected int
		response map[string]interface{}
	}

	tests := map[string]test{
		"simple": {
			input: map[string]string{
				"message": "value",
			},
			caller: func(store references.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(500))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Template: specs.Template{
							Message: specs.Message{
								"status": {
									Name:  "status",
									Path:  "status",
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
								"message": {
									Name:  "message",
									Path:  "message",
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
						},
					},
				},
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
							Default: "value",
						},
					},
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
			caller: func(store references.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(500))
				store.StoreValue("input", "status", int64(401))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Template: specs.Template{
							Message: specs.Message{
								"status": {
									Name:  "status",
									Path:  "status",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type: types.Int64,
										},
										Reference: &specs.PropertyReference{
											Resource: "input",
											Path:     "status",
										},
									},
								},
								"message": {
									Name:  "message",
									Path:  "message",
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
						},
					},
				},
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.Int64,
							Default: int64(401),
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.String,
							Default: "value",
						},
					},
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
			caller: func(store references.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(404))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Template: specs.Template{
							Message: specs.Message{
								"status": {
									Name:  "status",
									Path:  "status",
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
								"message": {
									Name:  "message",
									Path:  "message",
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
						},
					},
				},
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.Int64,
							Default: int64(404),
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.String,
							Default: "value",
						},
					},
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
			caller: func(store references.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(404))
			},
			err: &specs.OnError{
				Response: &specs.ParameterMap{
					Property: &specs.Property{
						Template: specs.Template{
							Message: specs.Message{
								"meta": {
									Name: "meta",
									Path: "meta",
									Template: specs.Template{
										Message: specs.Message{
											"status": {
												Name:  "status",
												Path:  "meta.status",
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
											"message": {
												Name:  "message",
												Path:  "meta.message",
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
									},
								},
								"const": {
									Name:  "const",
									Path:  "const",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type:    types.String,
											Default: "custom message",
										},
									},
								},
							},
						},
					},
				},
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.Int64,
							Default: int64(404),
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.String,
							Default: "value",
						},
					},
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
			caller: func(store references.Store) {
				store.StoreValue("error", "message", "value")
				store.StoreValue("error", "status", int64(404))
			},
			err: &specs.OnError{
				Response: nil,
				Status: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.Int64,
							Default: int64(404),
						},
					},
				},
				Message: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type:    types.String,
							Default: "value",
						},
					},
				},
			},

			expected: 404,
			response: nil,
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
