package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
)

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

	ctx := context.Background()
	req := protocol.Request{
		Context: ctx,
	}

	service := NewMockService(server.URL, "GET", "/")
	caller, err := (&Caller{}).New(service, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	defer caller.Close()

	r, w := io.Pipe()
	rw := &MockResponseWriter{
		header: protocol.Header{},
		writer: w,
	}

	go func() {
		caller.Call(rw, &req, refs)
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

func TestCallerUnkownMethod(t *testing.T) {
	service := NewMockService("http://localhost", "GET", "/")
	_, err := (&Caller{}).New(service, "unkown", nil)
	if err == nil {
		t.Fatal("unexpected pass expected a error to be thrown")
	}
}
