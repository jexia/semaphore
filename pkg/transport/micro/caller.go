package micro

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/codec/bytes"
	micrometa "github.com/micro/go-micro/v2/metadata"
	"go.uber.org/zap"
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
	return func(ctx *broker.Context) transport.Caller {
		return &Caller{
			name:    name,
			service: service,
		}
	}
}

// Caller represents the caller constructor
type Caller struct {
	ctx     *broker.Context
	name    string
	service Service
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return caller.name
}

// Dial constructs a new caller for the given service
func (caller *Caller) Dial(service *specs.Service, functions functions.Custom, opts specs.Options) (transport.Call, error) {
	module := broker.WithModule(caller.ctx, "caller", "go-micro")
	ctx := logger.WithFields(logger.WithLogger(module), zap.String("service", service.Name))

	methods := make(map[string]*Method, len(service.Methods))

	for _, method := range service.Methods {
		methods[method.Name] = &Method{
			name:       method.Name,
			endpoint:   fmt.Sprintf("%s.%s", service.Name, method.Name),
			references: make([]*specs.Property, 0),
		}
	}

	result := &Call{
		ctx:     ctx,
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
	ctx     *broker.Context
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
func (call *Call) SendMsg(ctx context.Context, rw transport.ResponseWriter, pr *transport.Request, refs references.Store) error {
	if pr.Method == nil {
		return ErrUndefinedMethod{}
	}

	bb, err := ioutil.ReadAll(pr.Body)
	if err != nil {
		return err
	}

	ctx = micrometa.NewContext(ctx, CopyMetadataHeader(pr.Header))

	method := call.methods[pr.Method.GetName()]
	if method == nil {
		return ErrUnknownMethod{
			Method: pr.Method.GetName(),
		}
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
			logger.Error(call.ctx, "unable to write the response body", zap.Error(err))
		}
	}()

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	logger.Info(call.ctx, "cloding go micro caller")
	return nil
}
