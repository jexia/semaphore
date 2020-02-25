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

	manager := new(Caller)
	caller := manager.Open(server.URL, nil)

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
