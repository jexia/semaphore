package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

type MockResponseWriter struct {
	header protocol.Header
	writer io.Writer
}

func (rw *MockResponseWriter) Header() protocol.Header {
	return rw.header
}

func (rw *MockResponseWriter) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

func (rw *MockResponseWriter) WriteHeader(int) {}

func AvailablePort(t *testing.T) int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func TestProxyForward(t *testing.T) {
	message := "hello world"
	specs := &specs.ParameterMap{
		Properties: map[string]*specs.Property{
			"message": &specs.Property{
				Name: "message",
				Path: "message",
				Type: types.TypeString,
			},
		},
	}

	codec, err := json.New("input", nil, specs)
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

	caller := NewCaller(server.URL, nil)

	r, w := io.Pipe()
	rw := &MockResponseWriter{
		header: protocol.Header{},
		writer: w,
	}

	go func() {
		caller.Call(rw, req)
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

func TestListener(t *testing.T) {
	message := "hello world"
	specs := &specs.ParameterMap{
		Properties: map[string]*specs.Property{
			"message": &specs.Property{
				Name: "message",
				Path: "message",
				Type: types.TypeString,
			},
		},
	}

	refs := refs.NewStore(1)
	codec, err := json.New("input", nil, specs)
	if err != nil {
		t.Fatal(err)
	}

	port := AvailablePort(t)
	addr := fmt.Sprintf(":%d", port)
	listener := NewListener(addr, nil)
	defer listener.Close()

	go listener.Serve(func(writer protocol.ResponseWriter, request protocol.Request) {
		codec.Unmarshal(request.Body, refs)
	})

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
	http.Post(endpoint, "application/json", strings.NewReader(`{"message":"`+message+`"}`))

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
