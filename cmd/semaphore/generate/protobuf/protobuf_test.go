package protobuf

import (
	"bytes"
	"context"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/codec/tests"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/jhump/protoreflect/desc/protoprint"
)

type flow struct{}

func (flow flow) NewStore() references.Store { return nil }

func (flow flow) GetName() string { return "Mock" }

func (flow flow) Errors() []transport.Error { return nil }

func (flow flow) Do(ctx context.Context, refs references.Store) error { return nil }

func (flow flow) Wait() {}

var (
	request = &specs.Property{
		Template: specs.Template{
			Message: make(specs.Message),
		},
	}

	response = &specs.Property{
		Name:  "root",
		Label: labels.Optional,
		Template: specs.Template{
			Identifier: "recursive",
			Message: specs.Message{
				"string": func() *specs.Property {
					var clone = tests.PropString()
					clone.Position = 1
					clone.Path = "root." + clone.Path

					return clone
				}(),
				"integer": func() *specs.Property {
					var clone = tests.PropInteger()
					clone.Position = 2
					clone.Path = "root." + clone.Path

					return clone
				}(),
			},
		},
	}

	recursive = &specs.Property{
		Name:     "meta",
		Path:     "meta",
		Position: 3,
		Label:    labels.Optional,
		Template: specs.Template{
			Identifier: "recursive",
			Message: specs.Message{
				"string": &specs.Property{
					Name:     "text",
					Path:     "meta.text",
					Position: 1,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			},
		},
	}

	expected = `syntax = "proto3";

package semaphore;

message ServiceMockRequest {
}

message ServiceMockResponse {
  string string = 1;

  int32 integer = 2;

  metaType meta = 3;

  message metaType {
    string text = 1;

    metaType meta = 3;
  }
}

service service {
  rpc Mock ( ServiceMockRequest ) returns ( ServiceMockResponse );
}
`
)

func init() {
	recursive.Message["meta"] = recursive
	response.Message["meta"] = recursive
}

func TestGenerate(t *testing.T) {
	var endpoints = transport.EndpointList{
		{
			Listener: "grpc",
			Flow:     flow{},
			Request: &transport.Object{
				Definition: &specs.ParameterMap{
					Property: request,
				},
			},
			Response: &transport.Object{
				Definition: &specs.ParameterMap{
					Property: response,
				},
			},
			Errs: nil,
		},
	}

	result, err := generate(broker.NewBackground(), endpoints)
	if err != nil {
		t.Fatal(err)
	}

	service, ok := result["semaphore.service"]
	if !ok {
		t.Fatal("result does not contain expected service")
	}

	descriptor, err := service.FileDescriptor()
	if err != nil {
		t.Fatal(err)
	}

	messages := descriptor.GetMessageTypes()
	if actual := len(messages); actual != 2 {
		t.Fatalf("result was expected to have 2 messages, got %d", actual)
	}

	var (
		buff    = bytes.NewBuffer([]byte{})
		printer = &protoprint.Printer{}
	)

	if err := printer.PrintProtoFile(descriptor, buff); err != nil {
		t.Fatal(err)
	}

	if actual := buff.String(); actual != expected {
		t.Errorf("unexpected output:\n%s", actual)
	}
}
