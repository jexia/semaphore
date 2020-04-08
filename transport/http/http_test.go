package http

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/jexia/maestro/internal/flow"
	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
)

type caller struct {
	fn func(context.Context, specs.Store) error
}

func (caller *caller) Do(ctx context.Context, store specs.Store) error {
	return caller.fn(ctx, store)
}

func (caller *caller) References() []*specs.Property {
	return nil
}

func NewCallerFunc(fn func(context.Context, specs.Store) error) flow.Call {
	return &caller{fn: fn}
}

func NewSimpleMockSpecs() *specs.ParameterMap {
	return &specs.ParameterMap{
		Property: &specs.Property{
			Type:  types.Message,
			Label: labels.Optional,
			Nested: map[string]*specs.Property{
				"message": {
					Name: "message",
					Path: "message",
					Type: types.String,
				},
			},
		},
	}
}

func NewMockService(host string, method string, endpoint string) *MockService {
	return &MockService{
		host: host,
		methods: []schema.Method{
			&MockMethod{
				name: "mock",
				options: schema.Options{
					MethodOption:   method,
					EndpointOption: endpoint,
				},
			},
		},
	}
}

type MockService struct {
	pkg           string
	name          string
	documentation string
	host          string
	codec         string
	transport     string
	methods       []schema.Method
	options       schema.Options
}

func (service *MockService) GetPackage() string {
	return service.pkg
}

func (service *MockService) GetFullyQualifiedName() string {
	return service.name
}

func (service *MockService) GetName() string {
	return service.name
}

func (service *MockService) GetComment() string {
	return service.documentation
}

func (service *MockService) GetHost() string {
	return service.host
}

func (service *MockService) GetCodec() string {
	return service.codec
}

func (service *MockService) GetTransport() string {
	return service.transport
}

func (service *MockService) GetMethod(name string) schema.Method {
	for _, method := range service.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

func (service *MockService) GetMethods() schema.Methods {
	return service.methods
}

func (service *MockService) GetOptions() schema.Options {
	return service.options
}

type MockMethod struct {
	name          string
	documentation string
	options       schema.Options
	input         schema.Property
	output        schema.Property
}

func (method *MockMethod) GetName() string {
	return method.name
}

func (method *MockMethod) GetComment() string {
	return method.documentation
}

func (method *MockMethod) GetInput() schema.Property {
	return method.input
}

func (method *MockMethod) GetOutput() schema.Property {
	return method.output
}

func (method *MockMethod) GetOptions() schema.Options {
	return method.options
}

type MockResponseWriter struct {
	header metadata.MD
	writer io.Writer
	status int
}

func (rw *MockResponseWriter) Header() metadata.MD {
	return rw.header
}

func (rw *MockResponseWriter) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

func (rw *MockResponseWriter) WriteHeader(status int) {
	rw.status = status
}

func AvailablePort(t *testing.T) int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
