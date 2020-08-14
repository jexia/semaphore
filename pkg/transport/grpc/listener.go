package grpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/jhump/protoreflect/desc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	rpcMeta "google.golang.org/grpc/metadata"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Listener represents a HTTP listener
type Listener struct {
	addr        string
	ctx         *broker.Context
	server      *grpc.Server
	methods     map[string]*Method
	descriptors map[string]*desc.FileDescriptor
	mutex       sync.RWMutex
}

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) transport.NewListener {
	return func(parent *broker.Context) transport.Listener {
		module := broker.WithModule(parent, "listener", "grpc")
		ctx := logger.WithLogger(module)

		return &Listener{
			addr: addr,
			ctx:  ctx,
		}
	}
}

// Name returns the name of the given listener
func (listener *Listener) Name() string { return "grpc" }

// Serve opens the HTTP listener and calls the given handler function on reach request
func (listener *Listener) Serve() error {
	logger.Info(listener.ctx, "serving gRPC listener", zap.String("addr", listener.addr))

	listener.server = grpc.NewServer(
		grpc.CustomCodec(Codec()),
		grpc.UnknownServiceHandler(listener.handler),
	)

	rpb.RegisterServerReflectionServer(listener.server, listener)

	lis, err := net.Listen("tcp", listener.addr)
	if err != nil {
		return err
	}

	return listener.server.Serve(lis)
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(ctx *broker.Context, endpoints []*transport.Endpoint, codecs map[string]codec.Constructor) error {
	logger.Info(listener.ctx, "gRPC listener received new endpoints")

	var (
		constructor = proto.NewConstructor()
		methods     = make(map[string]*Method, len(endpoints))
		services    = make(map[string]*proto.Service)
		descriptors = make(map[string]*desc.FileDescriptor)
	)

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint)
		if err != nil {
			return err
		}

		var (
			service = fmt.Sprintf("%s.%s", options.Package, options.Service)
			name    = fmt.Sprintf("%s/%s", service, options.Method)

			method = &Method{
				Endpoint: endpoint,
				fqn:      name,
				name:     options.Method,
				flow:     endpoint.Flow,
			}
		)

		if err := method.NewCodec(ctx, constructor); err != nil {
			return err
		}

		methods[name] = method

		if services[service] == nil {
			services[service] = &proto.Service{
				Package: options.Package,
				Name:    options.Service,
				Methods: make(proto.Methods),
			}
		}

		services[service].Methods[name] = methods[name]
	}

	for key, service := range services {
		descriptor, err := service.FileDescriptor()
		if err != nil {
			return fmt.Errorf("cannot generate file descriptor for %q: %w", key, err)
		}

		descriptors[key] = descriptor
	}

	listener.mutex.Lock()
	listener.methods = methods
	listener.descriptors = descriptors
	listener.mutex.Unlock()

	return nil
}

func (listener *Listener) handler(srv interface{}, stream grpc.ServerStream) error {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()

	fqn, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return grpc.Errorf(codes.Internal, "low level server stream not exists in context")
	}

	listener.mutex.RLock()
	method := listener.methods[fqn[1:]]
	listener.mutex.RUnlock()

	if method == nil {
		return grpc.Errorf(codes.Unimplemented, "unknown method: %s", fqn)
	}

	req := &frame{}
	err := stream.RecvMsg(req)
	if err != nil {
		return err
	}

	store := method.flow.NewStore()

	if method.Request != nil {
		header, ok := rpcMeta.FromIncomingContext(stream.Context())
		if ok {
			method.Request.Meta.Unmarshal(CopyRPCMD(header), store)
		}

		err = method.Request.Codec.Unmarshal(bytes.NewBuffer(req.payload), store)
		if err != nil {
			return grpc.Errorf(codes.ResourceExhausted, "invalid message body: %s", err)
		}
	}

	err = method.flow.Do(stream.Context(), store)
	if err != nil {
		object := method.Endpoint.Errs.Get(transport.Unwrap(err))
		if object == nil {
			logger.Error(listener.ctx, "unable to lookup error manager", zap.Error(err))
			return grpc.Errorf(codes.Internal, err.Error())
		}

		message := object.ResolveMessage(store)
		status := object.ResolveStatusCode(store)

		return grpc.Errorf(CodeFromStatus(status), message)
	}

	if method.Response != nil {
		header := method.Response.Meta.Marshal(store)
		reader, err := method.Response.Codec.Marshal(store)
		if err != nil {
			return grpc.Errorf(codes.ResourceExhausted, "invalid response body: %s", err)
		}

		bb, err := ioutil.ReadAll(reader)
		if err != nil {
			return grpc.Errorf(codes.ResourceExhausted, "unable to read full response body: %s", err)
		}

		res := &frame{
			payload: bb,
		}

		err = stream.SendMsg(res)
		if err != nil {
			return grpc.Errorf(codes.Internal, "unknown error: %s", err)
		}

		stream.SetTrailer(CopyMD(header))
	}

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	logger.Info(listener.ctx, "closing gRPC listener")
	listener.server.GracefulStop()
	return nil
}
