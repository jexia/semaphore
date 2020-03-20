package reflection

import (
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

func Handler(srv interface{}, stream grpc.ServerStream) error {
	log.Println("handler called!")
	srv.(rpb.ServerReflectionServer).ServerReflectionInfo(&serverReflectionServerReflectionInfoServer{stream})
	return nil
}

type serverReflectionServerReflectionInfoServer struct {
	grpc.ServerStream
}

func (x *serverReflectionServerReflectionInfoServer) Send(m *rpb.ServerReflectionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serverReflectionServerReflectionInfoServer) Recv() (*rpb.ServerReflectionRequest, error) {
	m := new(rpb.ServerReflectionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var desc = grpc.ServiceDesc{
	ServiceName: "grpc.reflection.v1alpha.ServerReflection",
	HandlerType: (*rpb.ServerReflectionServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerReflectionInfo",
			Handler:       Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc_reflection_v1alpha/reflection.proto",
}

func Register(server *grpc.Server) {
	server.RegisterService(&desc, &reflectionServer{server})
}

type reflectionServer struct {
	s *grpc.Server
}

// ServerReflectionInfo is the reflection service handler.
func (s *reflectionServer) ServerReflectionInfo(stream rpb.ServerReflection_ServerReflectionInfoServer) error {
	for {
		in, err := stream.Recv()
		log.Println("incoming!")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		out := &rpb.ServerReflectionResponse{
			ValidHost:       in.Host,
			OriginalRequest: in,
		}
		switch req := in.MessageRequest.(type) {
		case *rpb.ServerReflectionRequest_FileByFilename:
			log.Println("file name")
		case *rpb.ServerReflectionRequest_FileContainingSymbol:
			log.Println("file containing symbol")
		case *rpb.ServerReflectionRequest_FileContainingExtension:
			log.Println("file containing extension")
		case *rpb.ServerReflectionRequest_AllExtensionNumbersOfType:
			log.Println("file containing all extension")
		case *rpb.ServerReflectionRequest_ListServices:
			log.Println("file containing list services")

			out.MessageResponse = &rpb.ServerReflectionResponse_ListServicesResponse{
				ListServicesResponse: &rpb.ListServiceResponse{
					Service: []*rpb.ServiceResponse{
						{
							Name: "package.hello",
						},
					},
				},
			}
		default:
			log.Println(req)
			return status.Errorf(codes.InvalidArgument, "invalid MessageRequest: %v", in.MessageRequest)
		}

		if err := stream.Send(out); err != nil {
			return err
		}
	}
}
