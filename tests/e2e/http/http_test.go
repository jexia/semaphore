package http

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

	"github.com/jexia/semaphore/v2"
	"github.com/jexia/semaphore/v2/cmd/semaphore/daemon/providers"
	codecJSON "github.com/jexia/semaphore/v2/pkg/codec/json"
	codecXML "github.com/jexia/semaphore/v2/pkg/codec/xml"
	transportHTTP "github.com/jexia/semaphore/v2/pkg/transport/http"
	"github.com/jexia/semaphore/v2/tests/e2e"
)

const (
	SemaphorePort = 8080
	SemaphoreHost = "localhost"
)

var SemaphoreHTTPAddr = fmt.Sprintf("%s:%d", SemaphoreHost, SemaphorePort)

// EchoHandler creates an HTTP handler that returns the request body as a response.
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

// EchoRouter creates an HTTP router for testing.
func EchoRouter(t *testing.T) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/echo", EchoHandler(t))
	// TODO: add more handlers

	return mux
}

func TestHTTPTransport(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

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
				body := map[string]map[string]interface{}{
					"data": {
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

				expected := map[string]interface{}{
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
				body := map[string]map[string]interface{}{
					"data": {
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

				expected := map[string]interface{}{
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

				body := data{
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

				expected := map[string]interface{}{
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

				testServer := httptest.NewUnstartedServer(handler(t))
				testServer.Listener.Close()
				testServer.Listener = listener

				testServer.Start()
				defer testServer.Close()
			}

			var (
				config = e2e.Config{
					SemaphoreOptions: []semaphore.Option{
						semaphore.WithCodec(codecXML.NewConstructor()),
						semaphore.WithCodec(codecJSON.NewConstructor()),
						semaphore.WithCaller(transportHTTP.NewCaller()),
					},
					ProviderOptions: []providers.Option{
						providers.WithListener(transportHTTP.NewListener(fmt.Sprintf(":%d", SemaphorePort))),
					},
				}
				semaphore = e2e.Instance(t, test.flow, test.schema, config)
				path      = fmt.Sprintf("http://%s/%s", SemaphoreHTTPAddr, test.path)
				ctype     = fmt.Sprintf("application/%s", test.path)
				request   = bytes.NewBuffer(test.request)
			)

			defer semaphore.Close()

			ready, errs := semaphore.Serve()

			select {
			case <-ready:
			case err := <-errs:
				t.Fatalf("error happened: %s", err)
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
