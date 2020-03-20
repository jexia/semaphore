package grpc

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var proxyDesc = &grpc.StreamDesc{
	ServerStreams: false,
	ClientStreams: false,
}

// NewCaller constructs a new HTTP caller
func NewCaller() *Caller {
	return &Caller{}
}

// Caller represents the caller constructor
type Caller struct {
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return "grpc"
}

// Dial constructs a new caller for the given host
func (caller *Caller) Dial(service schema.Service, functions specs.CustomDefinedFunctions, opts schema.Options) (protocol.Call, error) {
	log.WithFields(log.Fields{
		"service": service.GetName(),
		"host":    service.GetHost(),
	}).Info("Constructing new gRPC caller")

	methods := make(map[string]*Method, len(service.GetMethods()))

	for _, method := range service.GetMethods() {
		methods[method.GetName()] = &Method{
			name:     method.GetName(),
			endpoint: fmt.Sprintf("/%s/%s", service.GetName(), method.GetName()),
		}
	}

	result := &Call{
		service: service.GetName(),
		host:    service.GetHost(),
		methods: methods,
	}

	return result, nil
}

// Method represents a call method
type Method struct {
	name     string
	endpoint string
}

// GetName returns the method name
func (method *Method) GetName() string {
	return method.name
}

// References returns the available references inside the given method
func (method *Method) References() []*specs.Property {
	return make([]*specs.Property, 0)
}

// Call represents the HTTP caller implementation
type Call struct {
	service string
	host    string
	methods map[string]*Method
}

// GetMethods returns the available methods within the given call
func (call *Call) GetMethods() []protocol.Method {
	result := make([]protocol.Method, 0, len(call.methods))

	for _, method := range call.methods {
		result = append(result, method)
	}

	return result
}

// GetMethod attempts to return a method matching the given name
func (call *Call) GetMethod(name string) protocol.Method {
	for _, method := range call.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// SendMsg opens a new connection to the configured host and attempts to send the given headers and stream
func (call *Call) SendMsg(ctx context.Context, rw protocol.ResponseWriter, req *protocol.Request, refs *refs.Store) error {
	log.WithFields(log.Fields{
		"service": call.service,
		"method":  req.Method,
	}).Debug("Calling gRPC caller")

	method := call.methods[req.Method.GetName()]
	if method == nil {
		return trace.New(trace.WithMessage("unkown method '%s' for service '%s'", req.Method, call.service))
	}

	conn, err := grpc.DialContext(ctx, call.host, grpc.WithCodec(Codec()), grpc.WithInsecure())
	if err != nil {
		return err
	}

	stream, err := conn.NewStream(ctx, proxyDesc, method.endpoint)
	if err != nil {
		return err
	}

	bb, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	stream.SendMsg(&frame{
		payload: bb,
	})

	// stream.Header()

	res := &frame{}
	err = stream.RecvMsg(res)
	if err != nil {
		return err
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	_, err = rw.Write(res.payload)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	log.WithField("host", call.host).Info("Closing gRPC caller")
	return nil
}
