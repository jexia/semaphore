package grpc

import (
	"io"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

// ServerReflectionInfo handles the gRPC reflection v1 alpha implementation.
func (listener *Listener) ServerReflectionInfo(stream rpb.ServerReflection_ServerReflectionInfoServer) error {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()

	for {
		in, err := stream.Recv()
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
			out.MessageResponse = &rpb.ServerReflectionResponse_FileDescriptorResponse{
				FileDescriptorResponse: &rpb.FileDescriptorResponse{FileDescriptorProto: [][]byte{}},
			}
		case *rpb.ServerReflectionRequest_FileContainingSymbol:
			descriptor, ok := listener.descriptors[req.FileContainingSymbol]
			if !ok {
				out.MessageResponse = &rpb.ServerReflectionResponse_ErrorResponse{
					ErrorResponse: &rpb.ErrorResponse{
						ErrorCode:    int32(codes.NotFound),
						ErrorMessage: "symbol not found",
					},
				}
				continue
			}

			bb, err := proto.Marshal(descriptor.AsFileDescriptorProto())
			if err != nil {
				out.MessageResponse = &rpb.ServerReflectionResponse_ErrorResponse{
					ErrorResponse: &rpb.ErrorResponse{
						ErrorCode:    int32(codes.Internal),
						ErrorMessage: err.Error(),
					},
				}
				continue
			}

			out.MessageResponse = &rpb.ServerReflectionResponse_FileDescriptorResponse{
				FileDescriptorResponse: &rpb.FileDescriptorResponse{FileDescriptorProto: [][]byte{bb}},
			}
		case *rpb.ServerReflectionRequest_ListServices:
			services := []*rpb.ServiceResponse{}

			for key := range listener.descriptors {
				services = append(services, &rpb.ServiceResponse{
					Name: key,
				})
			}

			out.MessageResponse = &rpb.ServerReflectionResponse_ListServicesResponse{
				ListServicesResponse: &rpb.ListServiceResponse{
					Service: services,
				},
			}
		default:
			return status.Errorf(codes.InvalidArgument, "invalid message request: %v", in.MessageRequest)
		}

		if err := stream.Send(out); err != nil {
			return err
		}
	}
}
