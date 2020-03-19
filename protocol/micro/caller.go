package micro

import (
	"io"
	"io/ioutil"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/micro/go-micro/v2/transport"
	log "github.com/sirupsen/logrus"
)

// New constructs a new go micro transport wrapper
func New(name string, transporter transport.Transport) *Caller {
	return &Caller{
		name:        name,
		transporter: transporter,
	}
}

// Caller represents the caller constructor
type Caller struct {
	name        string
	transporter transport.Transport
}

// Name returns the name of the given caller
func (caller *Caller) Name() string {
	return caller.name
}

// Dial constructs a new caller for the given host
func (caller *Caller) Dial(schema schema.Service, functions specs.CustomDefinedFunctions, opts schema.Options) (protocol.Call, error) {
	client, err := caller.transporter.Dial(schema.GetHost())
	if err != nil {
		return nil, err
	}

	methods := make(map[string]*Method, len(schema.GetMethods()))

	for _, method := range schema.GetMethods() {
		methods[method.GetName()] = &Method{
			name:       method.GetName(),
			references: make([]*specs.Property, 0),
		}
	}

	result := &Call{
		transport: client,
	}

	return result, nil
}

// Method represents a service method
type Method struct {
	name       string
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
	transport transport.Client
	methods   map[string]*Method
}

// GetMethods returns the available methods within the service caller
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

// SendMsg calls the configured host and attempts to call the given endpoint with the given headers and stream
func (call *Call) SendMsg(rw protocol.ResponseWriter, pr *protocol.Request, refs *refs.Store) error {
	bb, err := ioutil.ReadAll(pr.Body)
	if err != nil {
		return err
	}

	err = call.transport.Send(&transport.Message{
		Header: pr.Header,
		Body:   bb,
	})

	if err != nil && err != io.EOF {
		return err
	}

	res := &transport.Message{}
	err = call.transport.Recv(res)
	if err != nil {
		return err
	}

	for key, val := range res.Header {
		rw.Header().Set(key, val)
	}

	_, err = rw.Write(res.Body)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	log.Info("Closing go micro caller")
	return nil
}
