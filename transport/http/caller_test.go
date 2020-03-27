package http

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
	"github.com/jexia/maestro/transport"
)

func NewMockCaller() *Caller {
	ctx := instance.NewContext()
	caller := &Caller{
		ctx: ctx,
	}
	return caller
}

func TestCaller(t *testing.T) {
	message := "hello world"
	specs := NewSimpleMockSpecs()

	codec, err := (&json.Constructor{}).New("input", specs)
	if err != nil {
		t.Fatal(err)
	}

	refs := refs.NewStore(1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"` + message + `"}`))
	}))

	defer server.Close()

	service := NewMockService(server.URL, "GET", "/")
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

	go func() {
		caller.SendMsg(ctx, rw, &req, refs)
		w.Close()
	}()

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

	method := call.GetMethod("unkown")
	if method != nil {
		t.Fatal("unexpected method returned")
	}
}

func TestCallerReferences(t *testing.T) {
	expected := ":message"
	path := "message"
	resource := ".request"

	service := NewMockService("http://localhost", "GET", "/"+expected)
	call, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	method := call.GetMethod("mock")
	references := method.References()
	if len(references) != 1 {
		t.Fatalf("unexpected references %+v", references)
	}

	reference := references[0]
	if reference.Path != expected {
		t.Fatalf("unexpected path %s, expected %s", reference.Path, expected)
	}

	if reference.Reference.Resource != resource {
		t.Fatalf("unexpected reference resource %s, expected %s", reference.Reference.Resource, resource)
	}

	if reference.Reference.Path != path {
		t.Fatalf("unexpected reference path %s, expected %s", reference.Reference.Path, path)
	}
}

func TestCallerReferencesLookup(t *testing.T) {
	value := "1"
	expected := "/" + value

	path := "message"
	resource := ".request"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expected {
			t.Log("unexpected url", r.URL, " expected", expected)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	defer server.Close()

	service := NewMockService(server.URL, "GET", "/:message")
	caller, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	method := caller.GetMethod("mock")
	references := method.References()
	if len(references) != 1 {
		t.Fatalf("unexpected references %+v", references)
	}

	references[0].Type = types.String
	references[0].Label = labels.Optional

	store := refs.NewStore(1)
	ctx := context.Background()
	req := transport.Request{
		Method: method,
	}

	store.StoreValue(resource, path, value)

	rw := &MockResponseWriter{
		header: metadata.MD{},
		writer: ioutil.Discard,
	}

	err = caller.SendMsg(ctx, rw, &req, store)
	if err != nil {
		t.Fatal(err)
	}
}
