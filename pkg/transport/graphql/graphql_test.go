package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

func NewSimpleMockSpecs() *specs.ParameterMap {
	return &specs.ParameterMap{
		Property: &specs.Property{
			Template: specs.Template{
				Message: specs.Message{
					"first_name": {
						Name: "first_name",
						Path: "first_name",
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "John",
							},
						},
					},
					"last_name": {
						Name: "last_name",
						Path: "last_name",
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "Doe",
							},
						},
					},
				},
			},
		},
	}
}

func AvailableAddr(t *testing.T) int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func TestNewListener(t *testing.T) {
	type test struct {
		endpoints []*transport.Endpoint
		request   string
		expected  map[string]interface{}
	}

	ctx := logger.WithLogger(broker.NewBackground())

	tests := map[string]test{
		"no operation": {
			endpoints: []*transport.Endpoint{},
			request:   "",
			expected: map[string]interface{}{
				"data": nil,
				"errors": []interface{}{
					map[string]interface{}{
						"locations": []interface{}{},
						"message":   "Must provide an operation.",
					},
				},
			},
		},
		"single": {
			endpoints: []*transport.Endpoint{
				{
					Request: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
					Flow:    flow.NewManager(ctx, "test", nil, nil, nil, nil),
					Options: specs.Options{
						PathOption: "user",
					},
					Response: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
				},
			},
			request: "{user{first_name}}",
			expected: map[string]interface{}{
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"first_name": "John",
					},
				},
			},
		},
		"multiple fields": {
			endpoints: []*transport.Endpoint{
				{
					Request: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
					Flow:    flow.NewManager(ctx, "test", nil, nil, nil, nil),
					Options: specs.Options{
						PathOption: "user",
					},
					Response: transport.NewObject(NewSimpleMockSpecs(), nil, nil),
				},
			},
			request: "{user{first_name,last_name}}",
			expected: map[string]interface{}{
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"first_name": "John",
						"last_name":  "Doe",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			port := AvailableAddr(t)
			graphql := NewListener(fmt.Sprintf(":%d", port), nil)
			listener := graphql(ctx)

			go func() {
				err := listener.Serve()
				if err != nil {
					t.Error(err)
				}
			}()

			// Some CI pipelines take a little while before the listener is active
			time.Sleep(100 * time.Millisecond)

			defer listener.Close()

			err := listener.Handle(ctx, test.endpoints, nil)
			if err != nil {
				t.Fatal(err)
			}

			endpoint := fmt.Sprintf("http://127.0.0.1:%d/", port)
			req, err := json.Marshal(req{
				Query: test.request,
			})

			if err != nil {
				t.Fatal(err)
			}

			res, err := http.Post(endpoint, "application/json", bytes.NewReader(req))
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			body := map[string]interface{}{}
			err = json.Unmarshal(bb, &body)
			if err != nil {
				t.Fatal(err)
			}

			log.Println(body)
			if diff := deep.Equal(body, test.expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
