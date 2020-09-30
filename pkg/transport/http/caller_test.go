package http

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/references"
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

func NewMockCaller() *Caller {
	ctx := logger.WithLogger(broker.NewBackground())
	caller := &Caller{
		ctx: ctx,
	}
	return caller
}

func TestNewCaller(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	constructor := NewCaller()
	listener := constructor(ctx)
	if listener == nil {
		t.Fatal("unexpected result, expected a listener to be constructed")
	}

	if listener.Name() != "http" {
		t.Fatalf("unexpected name '%s', expected listener to be called http", listener.Name())
	}
}

func TestCaller(t *testing.T) {
	message := "hello world"
	mock := NewSimpleMockSpecs()

	codec, err := (&json.Constructor{}).New("input", mock)
	if err != nil {
		t.Fatal(err)
	}

	refs := references.NewReferenceStore(1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/Path?query=value"
		got := r.URL.String()
		if got != want {
			t.Errorf("expected server addressed by %s, but got %s", want, got)
		}
		w.Write([]byte(`{"message":"` + message + `"}`))
	}))

	defer server.Close()

	service := NewMockService(server.URL, "GET", "/Path?query=value")
	caller, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	defer caller.Close()

	r, w := io.Pipe()
	rw := &MockResponseWriter{
		header: metadata.MD{},
		writer: w,
	}

	ctx := context.Background()
	req := transport.Request{
		Method: caller.GetMethod("mock"),
	}

	caller.SendMsg(ctx, rw, &req, refs)

	err = codec.Unmarshal(r, refs)
	if err != nil {
		t.Fatal(err)
	}

	ref := refs.Load("input", "message")
	if ref == nil {
		t.Fatal("input:message reference not set")
	}

	result, is := ref.Value.(string)
	if !is {
		t.Fatal("input:message reference is not a string")
	}

	if result != message {
		t.Fatalf("unexpected input:message %s, expected %s", result, message)
	}
}

func TestCallerUnknownMethod(t *testing.T) {
	service := NewMockService("http://localhost", "GET", "/")
	call, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	method := call.GetMethod("unknown")
	if method != nil {
		t.Fatal("unexpected method returned")
	}

	if len(call.GetMethods()) != 1 {
		t.Fatalf("unexpected methods %+v, expected a single method to be defined", call.GetMethods())
	}
}

func TestErrUnknownMethod(t *testing.T) {
	type fields struct {
		Method  string
		Service string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Method: "get", Service: "getsources"},
			"unknown method 'get' for service 'getsources'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUnknownMethod{
				Method:  "get",
				Service: "getsources",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallerReferences(t *testing.T) {
	type reference struct {
		raw      string
		path     string
		resource string
	}

	type test struct {
		endpoint   string
		references []reference
	}

	tests := map[string]test{
		"message": {
			endpoint: "/:message",
			references: []reference{
				{
					raw:      ":message",
					path:     "message",
					resource: ".params",
				},
			},
		},
		"multiple": {
			endpoint: "/:message/:name",
			references: []reference{
				{
					raw:      ":message",
					path:     "message",
					resource: ".params",
				},
				{
					raw:      ":name",
					path:     "name",
					resource: ".params",
				},
			},
		},
		"nested": {
			endpoint: "/:message/:nested.name",
			references: []reference{
				{
					raw:      ":message",
					path:     "message",
					resource: ".params",
				},
				{
					raw:      ":nested.name",
					path:     "nested.name",
					resource: ".params",
				},
			},
		},
		"staticsuffix": {
			endpoint: "/:message/suffix",
			references: []reference{
				{
					raw:      ":message",
					path:     "message",
					resource: ".params",
				},
			},
		},
		"staticprefix": {
			endpoint: "/prefix/:message",
			references: []reference{
				{
					raw:      ":message",
					path:     "message",
					resource: ".params",
				},
			},
		},
		"static": {
			endpoint: "/prefix/:message/suffix",
			references: []reference{
				{
					raw:      ":message",
					path:     "message",
					resource: ".params",
				},
			},
		},
		"number": {
			endpoint: "/:message1",
			references: []reference{
				{
					raw:      ":message1",
					path:     "message1",
					resource: ".params",
				},
			},
		},
		"underscore": {
			endpoint: "/:message_one",
			references: []reference{
				{
					raw:      ":message_one",
					path:     "message_one",
					resource: ".params",
				},
			},
		},
		"hyphen": {
			endpoint: "/:message-one",
			references: []reference{
				{
					raw:      ":message-one",
					path:     "message-one",
					resource: ".params",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			service := NewMockService("http://localhost", "GET", test.endpoint)
			call, err := NewMockCaller().Dial(service, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			method := call.GetMethod("mock")
			references := method.References()
			if len(call.GetMethods()) != 1 {
				t.Fatalf("unexpected methods %+v, expected a single method to be defined", call.GetMethods())
			}

			if len(references) != len(test.references) {
				t.Fatalf("unexpected references '%+v', expected '%+v'", references, test.references)
			}

			for index, test := range test.references {
				reference := references[index]
				if reference.Path != test.raw {
					t.Fatalf("unexpected path %s, expected %s", reference.Path, test.raw)
				}

				if reference.Reference.Resource != test.resource {
					t.Fatalf("unexpected reference resource %s, expected %s", reference.Reference.Resource, test.resource)
				}

				if reference.Reference.Path != test.path {
					t.Fatalf("unexpected reference path %s, expected %s", reference.Reference.Path, test.path)
				}
			}
		})
	}
}

func TestCallerReferencesLookup(t *testing.T) {
	type test struct {
		endpoint string
		value    string
		call     string
		path     string
		resource string
	}

	tests := map[string]test{
		"message": {
			endpoint: "/:message",
			value:    "1",
			call:     "/1",
			path:     "message",
			resource: ".params",
		},
		"nested": {
			endpoint: "/:nested.message",
			value:    "value",
			call:     "/value",
			path:     "nested.message",
			resource: ".params",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != test.call {
					w.WriteHeader(http.StatusBadRequest)
				}
			}))

			defer server.Close()

			service := NewMockService(server.URL, "GET", test.endpoint)
			caller, err := NewMockCaller().Dial(service, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			method := caller.GetMethod("mock")
			refs := method.References()
			if len(refs) != 1 {
				t.Fatalf("unexpected references %+v", refs)
			}

			store := references.NewReferenceStore(1)
			ctx := context.Background()
			req := transport.Request{
				Method: method,
			}

			store.StoreValue(test.resource, test.path, test.value)

			rw := &MockResponseWriter{
				header: metadata.MD{},
				writer: &DiscardWriter{},
			}

			err = caller.SendMsg(ctx, rw, &req, store)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
