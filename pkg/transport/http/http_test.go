package http

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// JSONEqual compares the JSON from two Readers.
func JSONEqual(a, b io.Reader) (bool, interface{}, interface{}, error) {
	var j, j2 interface{}
	d := json.NewDecoder(a)
	if err := d.Decode(&j); err != nil {
		return false, j, j2, err
	}
	d = json.NewDecoder(b)
	if err := d.Decode(&j2); err != nil {
		return false, j, j2, err
	}
	return reflect.DeepEqual(j2, j), j, j2, nil
}

type caller struct {
	fn func(context.Context, references.Store) error
}

func (caller *caller) Do(ctx context.Context, store references.Store) error {
	return caller.fn(ctx, store)
}

func (caller *caller) References() []*specs.Property {
	return nil
}

func NewCallerFunc(fn func(context.Context, references.Store) error) flow.Call {
	return &caller{fn: fn}
}

func NewSimpleMockSpecs() *specs.ParameterMap {
	return &specs.ParameterMap{
		Header: specs.Header{
			"Authorization": &specs.Property{},
			"Timestamp":     &specs.Property{},
		},
		Property: &specs.Property{
			Template: specs.Template{
				Message: specs.Message{
					"message": {
						Name: "message",
						Path: "message",
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
	}
}

func NewMockService(host string, method string, endpoint string) *specs.Service {
	return &specs.Service{
		Host: host,
		Methods: []*specs.Method{
			{
				Name: "mock",
				Options: specs.Options{
					MethodOption:   method,
					EndpointOption: endpoint,
				},
			},
		},
	}
}

type MockResponseWriter struct {
	header  metadata.MD
	writer  io.WriteCloser
	status  int
	message string
}

func (rw *MockResponseWriter) Header() metadata.MD {
	return rw.header
}

func (rw *MockResponseWriter) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

func (rw *MockResponseWriter) HeaderStatus(status int) {
	rw.status = status
}

func (rw *MockResponseWriter) HeaderMessage(message string) {
	rw.message = message
}

func (rw *MockResponseWriter) Close() error {
	return rw.writer.Close()
}

func AvailablePort(t *testing.T) int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
