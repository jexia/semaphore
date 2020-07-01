package micro

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/core/trace"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/codec/bytes"
	micrometa "github.com/micro/go-micro/metadata"
)

// Service is an interface that wraps the lower level libraries
// within go-micro. Its a convenience method for building
// and initialising services.
type Service interface {
	// The service name
	Name() string
	// Client is used to call services
	Client() client.Client
	// The service implementation
	String() string
}

// NewCaller constructs a new go micro transport wrapper
func NewCaller(name string, service Service) transport.NewCaller {
	return func(ctx instance.Context) transport.Caller {
		return &Caller{
			name:    name,
			service: service,
		}
	}
}

// Caller represents the caller constructor
type Caller struct {
	ctx     instance.Context
	name    string
	service Service
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return caller.name
}

// Dial constructs a new caller for the given service
func (caller *Caller) Dial(service *specs.Service, functions functions.Custom, opts specs.Options) (transport.Call, error) {
	methods := make(map[string]*Method, len(service.Methods))

	for _, method := range service.Methods {
		methods[method.Name] = &Method{
			name:       method.Name,
			endpoint:   fmt.Sprintf("%s.%s", service.Name, method.Name),
			references: make([]*specs.Property, 0),
		}
	}

	result := &Call{
		ctx:     caller.ctx,
		pkg:     service.Package,
		service: service.Name,
		methods: methods,
		client:  caller.service.Client(),
	}

	return result, nil
}

// Method represents a service method
type Method struct {
	name       string
	endpoint   string
	references []*specs.Property
}

// GetName returns the method name
func (method *Method) GetName() string {
	return method.name
}

// References returns the available method references
func (method *Method) References() []*specs.Property {
	if method.references == nil {
		return make([]*specs.Property, 0)
	}

	return method.references
}

// Call represents the go micro transport wrapper implementation
type Call struct {
	ctx     instance.Context
	pkg     string
	service string
	client  client.Client
	methods map[string]*Method
}

// GetMethods returns the available methods within the service caller
func (call *Call) GetMethods() []transport.Method {
	result := make([]transport.Method, 0, len(call.methods))

	for _, method := range call.methods {
		result = append(result, method)
	}

	return result
}

// GetMethod attempts to return a method matching the given name
func (call *Call) GetMethod(name string) transport.Method {
	for _, method := range call.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// SendMsg calls the configured service and attempts to call the given endpoint with the given headers and stream
func (call *Call) SendMsg(ctx context.Context, rw transport.ResponseWriter, pr *transport.Request, refs refs.Store) error {
	if pr.Method == nil {
		return trace.New(trace.WithMessage("method required, proxy forward not supported"))
	}

	bb, err := ioutil.ReadAll(pr.Body)
	if err != nil {
		return err
	}

	ctx = micrometa.NewContext(ctx, CopyMetadataHeader(pr.Header))

	method := call.methods[pr.Method.GetName()]
	if method == nil {
		return trace.New(trace.WithMessage("unknown service method %s", pr.Method.GetName()))
	}

	req := call.client.NewRequest(call.pkg, method.endpoint, &bytes.Frame{
		Data: bb,
	})

	res := &bytes.Frame{
		Data: []byte{},
	}

	err = call.client.Call(ctx, req, res)
	if err != nil {
		return err
	}

	rw.HeaderStatus(transport.StatusOK)
	rw.HeaderMessage(transport.StatusMessage(transport.StatusOK))

	go func() {
		defer rw.Close()
		_, err = rw.Write(res.Data)
		if err != nil {
			call.ctx.Logger(logger.Transport).Error(err)
		}
	}()

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	call.ctx.Logger(logger.Transport).Info("Closing go micro caller")
	return nil
}
