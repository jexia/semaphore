package grpc

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/flow"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var (
	// fileDescriptor of each test proto file.
	fdTest       *dpb.FileDescriptorProto
	fdTestv3     *dpb.FileDescriptorProto
	fdProto2     *dpb.FileDescriptorProto
	fdProto2Ext  *dpb.FileDescriptorProto
	fdProto2Ext2 *dpb.FileDescriptorProto
	// fileDescriptor marshalled.
	fdTestByte       []byte
	fdTestv3Byte     []byte
	fdProto2Byte     []byte
	fdProto2ExtByte  []byte
	fdProto2Ext2Byte []byte
)

func NewMockServer(t *testing.T, endpoints []*transport.Endpoint) (*grpc.ClientConn, *Listener) {
	n, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	addr := n.Addr().String()
	n.Close()

	ctx := instance.NewContext()
	constructor := NewListener(addr, nil)
	listener := constructor(ctx).(*Listener)
	err = listener.Handle(ctx, endpoints, nil)
	if err != nil {
		t.Fatal(err)
	}

	go listener.Serve()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("cannot connect to server: %v", err)
	}

	return conn, listener
}

func TestListServices(t *testing.T) {
	ctx := instance.NewContext()
	endpoints := []*transport.Endpoint{
		{
			Options: specs.Options{
				PackageOption: "com.mock",
				ServiceOption: "first",
			},
			Flow: flow.NewManager(ctx, "Get", []*flow.Node{}, nil, nil, nil),
		},
		{
			Options: specs.Options{
				PackageOption: "com.mock",
				ServiceOption: "second",
			},
			Flow: flow.NewManager(ctx, "Get", []*flow.Node{}, nil, nil, nil),
		},
	}

	conn, _ := NewMockServer(t, endpoints)
	defer conn.Close()

	c := rpb.NewServerReflectionClient(conn)
	stream, err := c.ServerReflectionInfo(context.Background(), grpc.WaitForReady(true))
	if err != nil {
		t.Fatal(err)
	}

	if err := stream.Send(&rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_ListServices{},
	}); err != nil {
		t.Fatal(err)
	}

	r, err := stream.Recv()
	if err != nil {
		// io.EOF is not ok.
		t.Fatal(err)
	}

	switch r.MessageResponse.(type) {
	case *rpb.ServerReflectionResponse_ListServicesResponse:
		services := r.GetListServicesResponse().Service
		want := []string{
			"com.mock.first",
			"com.mock.second",
		}
		// Compare service names in response with want.
		if len(services) != len(want) {
			t.Errorf("= %v, want service names: %v", services, want)
		}
		m := make(map[string]int)
		for _, e := range services {
			m[e.Name]++
		}
		for _, e := range want {
			if m[e] > 0 {
				m[e]--
				continue
			}
			t.Errorf("ListService\nreceived: %v,\nwant: %q", services, want)
		}
	default:
		t.Errorf("ListServices = %v, want type <ServerReflectionResponse_ListServicesResponse>", r.MessageResponse)
	}
}

func TestFileContainingSymbol(t *testing.T) {
	ctx := instance.NewContext()
	endpoints := []*transport.Endpoint{
		{
			Options: specs.Options{
				PackageOption: "com.mock",
				ServiceOption: "first",
			},
			Flow: flow.NewManager(ctx, "Get", []*flow.Node{}, nil, nil, nil),
		},
		{
			Options: specs.Options{
				PackageOption: "com.mock",
				ServiceOption: "second",
			},
			Flow: flow.NewManager(ctx, "Get", []*flow.Node{}, nil, nil, nil),
		},
	}

	conn, listener := NewMockServer(t, endpoints)
	defer conn.Close()

	c := rpb.NewServerReflectionClient(conn)
	stream, err := c.ServerReflectionInfo(context.Background(), grpc.WaitForReady(true))
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		symbol string
		want   *Service
	}

	tests := map[string]test{
		"first": {
			symbol: "com.mock.first",
			want:   listener.services["com.mock.first"],
		},
		"second": {
			symbol: "com.mock.second",
			want:   listener.services["com.mock.second"],
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := stream.Send(&rpb.ServerReflectionRequest{
				MessageRequest: &rpb.ServerReflectionRequest_FileContainingSymbol{
					FileContainingSymbol: test.symbol,
				},
			})

			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			r, err := stream.Recv()
			if err != nil {
				// io.EOF is not ok.
				t.Fatalf("failed to recv response: %v", err)
			}

			expected, err := proto.Marshal(test.want.proto)
			if err != nil {
				t.Fatal(err)
			}

			switch r.MessageResponse.(type) {
			case *rpb.ServerReflectionResponse_FileDescriptorResponse:
				if !reflect.DeepEqual(r.GetFileDescriptorResponse().FileDescriptorProto[0], expected) {
					t.Errorf("FileContainingSymbol(%v)\nreceived: %q,\nwant: %q", test.symbol, r.GetFileDescriptorResponse().FileDescriptorProto[0], expected)
				}
			default:
				t.Errorf("FileContainingSymbol(%v) = %v, want type <ServerReflectionResponse_FileDescriptorResponse>", test.symbol, r.MessageResponse)
			}
		})
	}
}
