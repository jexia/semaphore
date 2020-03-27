package grpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) transport.NewListener {
	// options, err := ParseListenerOptions(opts)
	// if err != nil {
	// 	// TODO: log err
	// }

	return func(ctx instance.Context) transport.Listener {
		return &Listener{
			addr: addr,
			ctx:  ctx,
		}
	}
}

// Listener represents a HTTP listener
type Listener struct {
	addr      string
	ctx       instance.Context
	server    *grpc.Server
	endpoints map[string]*Endpoint
	mutex     sync.RWMutex
}

// Name returns the name of the given listener
func (listener *Listener) Name() string {
	return "grpc"
}

// Serve opens the HTTP listener and calls the given handler function on reach request
func (listener *Listener) Serve() error {
	listener.ctx.Logger(logger.Transport).WithField("addr", listener.addr).Info("Serving gRPC listener")

	listener.server = grpc.NewServer(
		grpc.CustomCodec(Codec()),
		grpc.UnknownServiceHandler(listener.handler),
	)

	lis, err := net.Listen("tcp", listener.addr)
	if err != nil {
		return err
	}

	err = listener.server.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(endpoints []*transport.Endpoint, codecs map[string]codec.Constructor) error {
	logger := listener.ctx.Logger(logger.Transport)
	logger.Info("gRPC listener received new endpoints")

	constructor := proto.NewConstructor()
	result := make(map[string]*Endpoint, len(endpoints))

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint)
		if err != nil {
			return err
		}

		in, err := constructor.New(specs.InputResource, endpoint.Request)
		if err != nil {
			return err
		}

		out, err := constructor.New(specs.OutputResource, endpoint.Response)
		if err != nil {
			return err
		}

		method := fmt.Sprintf("/%s.%s/%s", options.Package, options.Service, options.Method)
		result[method] = &Endpoint{
			name: method,
			flow: endpoint.Flow,
			in:   in,
			out:  out,
		}
	}

	listener.mutex.Lock()
	listener.endpoints = result
	listener.mutex.Unlock()

	return nil
}

func (listener *Listener) handler(_ interface{}, stream grpc.ServerStream) error {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()

	method, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return grpc.Errorf(codes.Internal, "low level server stream not exists in context")
	}

	_, ok = metadata.FromIncomingContext(stream.Context())
	if ok {
		// TODO: support header values
	}

	endpoint := listener.endpoints[method]
	if endpoint == nil {
		return grpc.Errorf(codes.Unimplemented, "unknown method: %s", method)
	}

	req := &frame{}
	err := stream.RecvMsg(req)
	if err != nil {
		return err
	}

	store := endpoint.flow.NewStore()
	err = endpoint.in.Unmarshal(bytes.NewBuffer(req.payload), store)
	if err != nil {
		return grpc.Errorf(codes.ResourceExhausted, "invalid message body: %s", err)
	}

	err = endpoint.flow.Call(stream.Context(), store)
	if err != nil {
		return grpc.Errorf(codes.Internal, "unkown error: %s", err)
	}

	reader, err := endpoint.out.Marshal(store)
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
		return grpc.Errorf(codes.Internal, "unkown error: %s", err)
	}

	stream.SetTrailer(metadata.MD{})

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	listener.ctx.Logger(logger.Transport).Info("Closing gRPC listener")
	listener.server.GracefulStop()
	return nil
}

// Endpoint represents a gRPC endpoint
type Endpoint struct {
	name string
	flow transport.Flow
	in   codec.Manager
	out  codec.Manager
}
