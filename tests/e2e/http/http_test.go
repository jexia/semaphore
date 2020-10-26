package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	codecJSON "github.com/jexia/semaphore/pkg/codec/json"
	codecProto "github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	transport "github.com/jexia/semaphore/pkg/transport/http"
)

func EchoHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				t.Errorf("failed to close request body: %s", err)
			}
		}()

		if _, err := io.Copy(w, r.Body); err != nil {
			t.Errorf("failed to send the reply: %s", err)
		}
	}
}

func EchoRouter(t *testing.T) http.Handler {
	var mux = http.NewServeMux()

	mux.Handle("/echo", EchoHandler(t))
	// TODO: add more handlers

	return mux
}

func Semaphore(t *testing.T, flow, schema string) *daemon.Client {
	ctx := logger.WithLogger(broker.NewContext())

	core, err := semaphore.NewOptions(ctx,
		semaphore.WithLogLevel("*", "error"),
		semaphore.WithFlows(hcl.FlowsResolver(flow)),
		semaphore.WithCodec(codecJSON.NewConstructor()),
		semaphore.WithCodec(codecProto.NewConstructor()),
		semaphore.WithCaller(transport.NewCaller()),
	)

	if err != nil {
		t.Fatalf("cannot instantiate semaphore core: %s", err)
	}

	options, err := providers.NewOptions(ctx, core,
		providers.WithEndpoints(hcl.EndpointsResolver(flow)),
		providers.WithSchema(protobuffers.SchemaResolver([]string{"./proto"}, schema)),
		providers.WithServices(protobuffers.ServiceResolver([]string{"./proto"}, schema)),
		providers.WithListener(transport.NewListener(":8080")),
	)

	if err != nil {
		t.Fatalf("unable to configure provider options: %s", err)
	}

	client, err := daemon.NewClient(ctx, core, options)
	if err != nil {
		t.Fatalf("failed to create a semaphore instance: %s", err)
	}

	return client
}

func TestSemaphore(t *testing.T) {
	type test struct {
		title     string
		flow      string
		schema    string
		resources map[string]func(*testing.T) http.Handler
		request   interface{}
		status    int
		response  interface{}
	}

	tests := []test{
		{
			title:  "echo",
			flow:   "./flow/echo.hcl",
			schema: "./proto/echo.proto",
			request: map[string]map[string]interface{}{
				"data": map[string]interface{}{
					"enum":    "ON",
					"string":  "foo",
					"integer": 42,
					"double":  3.14159,
					"numbers": []float64{1, 2, 3, 4, 5},
				},
			},
			status: http.StatusOK,
			response: map[string]interface{}{
				"echo": map[string]interface{}{
					"enum":    "ON",
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
			title:  "echo with intermediate resource",
			flow:   "./flow/echo_intermediate.hcl",
			schema: "./proto/echo.proto",
			resources: map[string]func(t *testing.T) http.Handler{
				":8081": EchoRouter,
			},
			request: map[string]map[string]interface{}{
				"data": map[string]interface{}{
					"enum":    "ON",
					"string":  "foo",
					"integer": 42,
					"double":  3.14159,
					"numbers": []float64{1, 2, 3, 4, 5},
				},
			},
			status: http.StatusOK,
			response: map[string]interface{}{
				"echo": map[string]interface{}{
					"enum":    "ON",
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
				defer listener.Close()

				var testServer = httptest.NewUnstartedServer(handler(t))
				testServer.Listener.Close()
				testServer.Listener = listener

				testServer.Start()
				defer testServer.Close()
			}

			var semaphore = Semaphore(t, test.flow, test.schema)
			defer semaphore.Close()

			ready, errs := semaphore.Serve()

			select {
			case <-ready:
			case err := <-errs:
				t.Fatalf("error happened: %s", err)
			}

			req, err := json.Marshal(test.request)
			if err != nil {
				t.Fatalf("unable to marshal the request: %s", err)
			}

			res, err := http.Post(fmt.Sprintf("%s/json", "http://localhost:8080"), "application/json", bytes.NewBuffer(req))
			if err != nil {
				t.Fatal(err)
			}

			if actual := res.StatusCode; actual != test.status {
				t.Errorf("got status [%d] was expected to be [%d]", actual, test.status)
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("cannot read the response body: %s", err)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatalf("failed to unmarshal the response: %s", err)
			}

			if !reflect.DeepEqual(response, test.response) {
				t.Errorf("the output\n[%+v]\n was expected to be\n[%+v]", response, test.response)
			}
		})
	}
}
