package grpc

import (
	"context"
	"net"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	"github.com/mwitkow/grpc-proxy/proxy"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) protocol.Listener {
	server := grpc.NewServer(grpc.CustomCodec(proxy.Codec()))
	return &Listener{
		server: server,
	}
}

// Listener represents a GraphQL listener
type Listener struct {
	server *grpc.Server
}

// Name returns the name of the given listener
func (listener *Listener) Name() string {
	return "grpc"
}

// Serve opens the GraphQL listener and calls the given handler function on reach request
func (listener *Listener) Serve() error {
	log.Println("serving! 50051")
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	log.Println("listner!")
	return listener.server.Serve(l)
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(endpoints []*protocol.Endpoint, constructors map[string]codec.Constructor) error {
	services := map[string][]string{}

	for _, endpoint := range endpoints {
		options, err := ParseListenerOptions(endpoint.Options)
		if err != nil {
			return err
		}

		if options.Service == "" {
			// NOTE: throw a error or log unkown service?
			continue
		}

		method := options.Method
		if method == "" {
			method = endpoint.Flow.GetName()
		}

		services[options.Service] = append(services[options.Service], method)
	}

	// reflection.Register(listener.server)

	for service, methods := range services {
		desc := ServiceDescriptor(service, methods)
		listener.server.RegisterService(desc, proxy.TransparentHandler(nil))
	}

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	log.Info("Closing gRPC listener")
	listener.server.GracefulStop()
	return nil
}

func ServiceDescriptor(service string, methods []string) *grpc.ServiceDesc {
	log.Println("service descriptor", service, methods)
	desc := &grpc.ServiceDesc{
		ServiceName: service,
		HandlerType: (*interface{})(nil),
	}

	for _, name := range methods {
		impl := grpc.MethodDesc{
			MethodName: name,
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				log.Println("handler called!")
				return nil, nil
			},
		}

		desc.Methods = append(desc.Methods, impl)
	}

	return desc
}
