package e2e

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	"github.com/jexia/semaphore/cmd/semaphore/middleware"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	codecJSON "github.com/jexia/semaphore/pkg/codec/json"
	codecProto "github.com/jexia/semaphore/pkg/codec/proto"
	codecXML "github.com/jexia/semaphore/pkg/codec/xml"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	transport "github.com/jexia/semaphore/pkg/transport/http"
)

const (
	semaphorePort = 8080
	semaphoreHost = "http://localhost"
)

var semaphoreURL = fmt.Sprintf("%s:%d", semaphoreHost, semaphorePort)

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
		semaphore.WithCodec(codecXML.NewConstructor()),
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
		providers.WithListener(transport.NewListener(fmt.Sprintf(":%d", semaphorePort))),
		providers.WithAfterConstructor(middleware.ServiceSelector(flow)),
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
		disabled  bool
		title     string
		flow      string
		schema    string
		resources map[string]func(*testing.T) http.Handler
		path      string
		request   []byte
		status    int
		assert    func(t *testing.T, data []byte)
	}

	tests := []test{
		{
			title:  "JSON echo",
			flow:   "./flow/echo.hcl",
			schema: "./proto/echo.proto",
			path:   "json",
			request: func(t *testing.T) []byte {
				var body = map[string]map[string]interface{}{
					"data": map[string]interface{}{
						"enum":    "ON",
						"string":  "foo",
						"integer": 42,
						"double":  3.14159,
						"numbers": []float64{1, 2, 3, 4, 5},
						// TODO: send recursive types
					},
				}

				encoded, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("unable to marshal the request: %s", err)
				}

				return encoded
			}(t),
			status: http.StatusOK,
			assert: func(t *testing.T, data []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(data, &response); err != nil {
					t.Fatalf("failed to unmarshal the response: %s", err)
				}

				var expected = map[string]interface{}{
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
				}

				if !reflect.DeepEqual(response, expected) {
					t.Errorf("the output\n[%+v]\n was expected to be\n[%+v]", response, expected)
				}
			},
		},
		{
			title:  "JSON echo with intermediate resource",
			flow:   "./flow/echo_intermediate.hcl",
			schema: "./proto/echo.proto",
			resources: map[string]func(t *testing.T) http.Handler{
				":8081": EchoRouter,
			},
			path: "json",
			request: func(t *testing.T) []byte {
				var body = map[string]map[string]interface{}{
					"data": map[string]interface{}{
						"enum":    "ON",
						"string":  "foo",
						"integer": 42,
						"double":  3.14159,
						"numbers": []float64{1, 2, 3, 4, 5},
					},
				}

				encoded, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("unable to marshal the request: %s", err)
				}

				return encoded
			}(t),
			status: http.StatusOK,
			assert: func(t *testing.T, data []byte) {
				var response map[string]interface{}
				if err := json.Unmarshal(data, &response); err != nil {
					t.Fatalf("failed to unmarshal the response: %s", err)
				}

				var expected = map[string]interface{}{
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
				}

				if !reflect.DeepEqual(response, expected) {
					t.Errorf("the output\n[%+v]\n was expected to be\n[%+v]", response, expected)
				}
			},
		},
		{
			disabled: true, // disabled until XML codec is fixed
			title:    "XML echo",
			flow:     "./flow/echo.hcl",
			schema:   "./proto/echo.proto",
			path:     "xml",
			request: func(t *testing.T) []byte {
				type request struct {
					Enum    string  `xml:"enum"`
					String  string  `xml:"string"`
					Integer int     `xml:"integer"`
					Float   float64 `xml:"double"`
					Numbers []int   `xml:"numbers"`
				}

				type data struct {
					Data request `xml:"data"`
				}

				var body = data{
					Data: request{
						Enum:    "ON",
						String:  "foo",
						Integer: 42,
						Float:   3.14159,
						Numbers: []int{1, 2, 3, 4, 5},
						// TODO: check recursive types
					},
				}

				encoded, err := xml.Marshal(body)
				if err != nil {
					t.Fatalf("unable to marshal the request: %s", err)
				}

				return encoded
			}(t),
			status: http.StatusOK,
			assert: func(t *testing.T, data []byte) {
				var response map[string]interface{}
				if err := xml.Unmarshal(data, &response); err != nil {
					t.Fatalf("failed to unmarshal the response: %s", err)
				}

				var expected = map[string]interface{}{
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
				}

				if !reflect.DeepEqual(response, expected) {
					t.Errorf("the output\n[%+v]\n was expected to be\n[%+v]", response, expected)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.disabled {
				t.Skip()
			}

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

			var (
				semaphore = Semaphore(t, test.flow, test.schema)
				path      = fmt.Sprintf("%s/%s", semaphoreURL, test.path)
				ctype     = fmt.Sprintf("application/%s", test.path)
				request   = bytes.NewBuffer(test.request)
			)

			defer semaphore.Close()

			ready, errs := semaphore.Serve()

			select {
			case <-ready:
			case err := <-errs:
				t.Logf("error happened: %s", err)
			}

			res, err := http.Post(path, ctype, request)
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

			test.assert(t, body)
		})
	}
}
