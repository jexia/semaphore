package protobuf

import (
	"bytes"
	"context"
	"strings"
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
		Template: &specs.Template{
			Message: make(specs.Message),
		},
	}

	response = &specs.Property{
		Name:  "root",
		Label: labels.Optional,
		Template: &specs.Template{
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
		Template: &specs.Template{
			Identifier: "recursive",
			Message: specs.Message{
				"string": &specs.Property{
					Name:     "string",
					Path:     "meta.string",
					Position: 1,
					Template: &specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			},
		},
	}

	// since we are not sure about the order of fields provided by generator
	// we only can check if all expected lines are there
	expected = map[string]int{
		`syntax = "proto3"`:           1,
		`package semaphore;`:          1,
		`message ServiceMockRequest`:  1,
		`message ServiceMockResponse`: 1,
		`string string = 1;`:          2,
		`int32 integer = 2;`:          1,
		`metaType meta = 3;`:          2,
		`message metaType`:            1,
		`service service`:             1,
		`rpc Mock ( ServiceMockRequest ) returns ( ServiceMockResponse );`: 1,
	}
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

	actual := buff.String()

	for substr, occurs := range expected {
		if total := count(actual, substr); total != occurs {
			t.Errorf("was expected to meet a substring %q x %d times", substr, occurs)
		}
	}
}

func count(input, substr string) int {
	var index int

	index = strings.Index(input, substr)

	if index < 0 {
		return 0
	}

	return 1 + count(input[index+1:], substr)
}
