package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/jexia/semaphore/tests/e2e"
	proto "github.com/jexia/semaphore/tests/e2e/grpc/proto"
	"google.golang.org/grpc"
)

var SemaphoreGRPCAddr = fmt.Sprintf("%s:%d", e2e.SemaphoreHost, e2e.SemaphoreGRPCPort)

// convert any interface{} to proto struct and back (use with recursive types).
func convert(t *testing.T, src, dst interface{}) {
	encoded, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("unable to encode: %s", err)
	}

	if err := json.Unmarshal(encoded, dst); err != nil {
		t.Fatalf("unable to encode: %s", err)
	}
}

type external struct {
	*proto.UnimplementedExternalServer
}

func (*external) Post(_ context.Context, input *proto.Data) (*proto.Data, error) {
	return input, nil
}

func TestGRPCTransport(t *testing.T) {
	type test struct {
		title     string
		flow      string
		schema    string
		resources map[string]func(*grpc.Server)
		request   interface{}
		expected  interface{}
	}

	tests := []test{
		{
			title:  "PROTO echo",
			flow:   "./flow/echo.hcl",
			schema: "./proto/echo.proto",
			request: map[string]interface{}{
				"data": map[string]interface{}{
					"enum":    1,
					"string":  "foo",
					"integer": 42,
					"double":  3.14159,
					"numbers": []int{1, 2, 3, 4, 5},
				},
			},
			expected: map[string]interface{}{
				"echo": map[string]interface{}{
					"enum":    float64(1),
					"string":  "foo",
					"integer": float64(42),
					"double":  float64(3.14159),
					"numbers": []interface{}{
						float64(1),
						float64(2),
						float64(3),
						float64(4),
						float64(5),
					},
				},
			},
		},
		{
			title:  "PROTO echo with intermediate resource",
			flow:   "./flow/echo_intermediate.hcl",
			schema: "./proto/echo.proto",
			resources: map[string]func(*grpc.Server){
				":5052": func(server *grpc.Server) {
					proto.RegisterExternalServer(server, new(external))
				},
			},
			request: map[string]interface{}{
				"data": map[string]interface{}{
					"enum":    1,
					"string":  "foo",
					"integer": 42,
					"double":  3.14159,
					"numbers": []int{1, 2, 3, 4, 5},
				},
			},
			expected: map[string]interface{}{
				"echo": map[string]interface{}{
					"enum":    float64(1),
					"string":  "foo",
					"integer": float64(42),
					"double":  float64(3.14159),
					"numbers": []interface{}{
						float64(1),
						float64(2),
						float64(3),
						float64(4),
						float64(5),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			for addr, handler := range test.resources {
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					t.Fatalf("unable to create a listener: %s", err)
				}
				defer func(listener io.Closer) {
					listener.Close()
				}(listener)

				server := grpc.NewServer()
				defer func(server interface{ Stop() }) {
					server.Stop()
				}(server)

				handler(server)

				go func(server *grpc.Server) {
					if err := server.Serve(listener); err != nil {
						t.Logf("failed to serve: %s", err)
					}
				}(server)

			}

			var semaphore = e2e.Instance(t, test.flow, test.schema)
			defer semaphore.Close()

			ready, errs := semaphore.Serve()

			select {
			case <-ready:
			case err := <-errs:
				t.Fatalf("error happened: %s", err)
			}

			// Set up a connection to the server.
			conn, err := grpc.Dial(SemaphoreGRPCAddr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				t.Fatalf("cannot connect to semaphore: %s", err)
			}
			defer conn.Close()

			var (
				client      = proto.NewTypetestClient(conn)
				request     = new(proto.Request)
				ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
			)

			defer cancel()

			convert(t, test.request, request)

			response, err := client.Run(ctx, request)
			if err != nil {
				t.Fatalf("could not perform the request: %s", err)
			}

			var output interface{}

			convert(t, response, &output)

			if !reflect.DeepEqual(output, test.expected) {
				t.Errorf("the output\n[%+v]\n was expected to be\n[%+v]", output, test.expected)
			}
		})
	}
}
