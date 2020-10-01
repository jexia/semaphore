package transport

import (
	"context"
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// MockFlowManager is used to mock a flow manager
type MockFlowManager struct {
	errs []Error
}

func (mock *MockFlowManager) NewStore() references.Store                          { return nil }
func (mock *MockFlowManager) GetName() string                                     { return "" }
func (mock *MockFlowManager) Errors() []Error                                     { return mock.errs }
func (mock *MockFlowManager) Do(ctx context.Context, refs references.Store) error { return nil }
func (mock *MockFlowManager) Wait()                                               {}

type MockCodecConstructor struct {
	err error
}

func (mock *MockCodecConstructor) Name() string { return "mock" }
func (mock *MockCodecConstructor) New(resource string, specs *specs.ParameterMap) (codec.Manager, error) {
	return nil, mock.err
}

func TestNewEndpointNil(t *testing.T) {
	endpoint := NewEndpoint("", nil, nil, nil, nil, nil)
	if endpoint == nil {
		t.Fatal("unexpected empty endpoint")
	}
}

func TestNewEndpoint(t *testing.T) {
	endpoint := NewEndpoint("http", &MockFlowManager{}, &Forward{}, specs.Options{}, &specs.ParameterMap{}, &specs.ParameterMap{})
	if endpoint == nil {
		t.Fatal("unexpected empty endpoint")
	}
}

func TestNewObjectNil(t *testing.T) {
	object := NewObject(nil, nil, nil)
	if object == nil {
		t.Fatal("unexpected empty object")
	}
}

func TestNewObject(t *testing.T) {
	object := NewObject(&specs.ParameterMap{}, &specs.Property{}, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}
}

func TestResolveStatusCode(t *testing.T) {
	expected := int64(200)
	status := &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.Int64,
			},
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "status",
			},
		},
	}

	object := NewObject(&specs.ParameterMap{}, status, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	store := references.NewReferenceStore(1)
	store.StoreValue("input", "status", expected)

	result := object.ResolveStatusCode(store)
	if int64(result) != expected {
		t.Errorf("unexpected result %d, expected %d", result, expected)
	}
}

func TestResolveStatusCodeNil(t *testing.T) {
	expected := StatusInternalErr
	object := NewObject(&specs.ParameterMap{}, nil, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	store := references.NewReferenceStore(0)
	result := object.ResolveStatusCode(store)
	if result != expected {
		t.Errorf("unexpected result %d, expected %d", result, expected)
	}
}

func TestResolveStatusCodeNilReference(t *testing.T) {
	expected := StatusInternalErr
	status := &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.Int64,
			},
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "status",
			},
		},
	}

	object := NewObject(&specs.ParameterMap{}, status, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	store := references.NewReferenceStore(0)
	result := object.ResolveStatusCode(store)
	if result != expected {
		t.Errorf("unexpected result %d, expected %d", result, expected)
	}
}
func TestResolveStatusMessage(t *testing.T) {
	expected := "unexpected mock err"
	message := &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.String,
			},
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "status",
			},
		},
	}

	object := NewObject(&specs.ParameterMap{}, &specs.Property{}, message)
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	store := references.NewReferenceStore(1)
	store.StoreValue("input", "status", expected)

	result := object.ResolveMessage(store)
	if result != expected {
		t.Errorf("unexpected result %s, expected %s", result, expected)
	}
}

func TestResolveStatusMessageNil(t *testing.T) {
	expected := StatusMessage(StatusInternalErr)
	object := NewObject(&specs.ParameterMap{}, &specs.Property{}, nil)
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	store := references.NewReferenceStore(0)
	result := object.ResolveMessage(store)
	if result != expected {
		t.Errorf("unexpected result %s, expected %s", result, expected)
	}
}

func TestResolveStatusMessageNilReference(t *testing.T) {
	expected := StatusMessage(StatusInternalErr)
	message := &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.Int64,
			},
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "status",
			},
		},
	}

	object := NewObject(&specs.ParameterMap{}, &specs.Property{}, message)
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	store := references.NewReferenceStore(0)
	result := object.ResolveMessage(store)
	if result != expected {
		t.Errorf("unexpected result %s, expected %s", result, expected)
	}
}

func TestObjectNewMeta(t *testing.T) {
	schema := &specs.ParameterMap{
		Header: specs.Header{
			"key": &specs.Property{
				Label: labels.Optional,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.String,
					},
				},
			},
		},
	}

	object := NewObject(schema, &specs.Property{}, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	ctx := logger.WithLogger(broker.NewBackground())
	object.NewMeta(ctx, "mock")

	if object.Meta == nil {
		t.Fatal("meta manager not set")
	}
}

func TestObjectNewMetaNil(t *testing.T) {
	object := NewObject(nil, &specs.Property{}, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	ctx := logger.WithLogger(broker.NewBackground())
	object.NewMeta(ctx, "mock")

	if object.Meta != nil {
		t.Fatal("unexpected meta manager set")
	}
}

func TestNewCodec(t *testing.T) {
	schema := &specs.ParameterMap{
		Property: &specs.Property{
			Template: specs.Template{
				Message: specs.Message{
					"key": {
						Name:  "key",
						Path:  "key",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
	}

	object := NewObject(schema, &specs.Property{}, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	ctx := logger.WithLogger(broker.NewBackground())
	err := object.NewCodec(ctx, "mock", json.NewConstructor())
	if err != nil {
		t.Fatal(err)
	}

	if object.Codec == nil {
		t.Fatal("codec manager not set")
	}
}

func TestNewCodecNil(t *testing.T) {
	object := NewObject(nil, &specs.Property{}, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	ctx := logger.WithLogger(broker.NewBackground())
	err := object.NewCodec(ctx, "mock", json.NewConstructor())
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCodecErr(t *testing.T) {
	expected := errors.New("mock error")
	mock := &MockCodecConstructor{err: expected}

	schema := &specs.ParameterMap{
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		},
	}

	object := NewObject(schema, &specs.Property{}, &specs.Property{})
	if object == nil {
		t.Fatal("unexpected empty object")
	}

	ctx := logger.WithLogger(broker.NewBackground())
	err := object.NewCodec(ctx, "mock", mock)
	if err != expected {
		t.Errorf("unexpected result %+v, expected %+v", err, expected)
	}
}

func TestSetErrs(t *testing.T) {
	typed := errors.New("mock err")
	key := WrapError(typed, &specs.OnError{})
	expected := NewObject(nil, nil, nil)

	collection := Errs{}
	collection.Set(key, expected)

	result := collection.Get(key)
	if result != expected {
		t.Errorf("unexpected result %+v, expected %+v", result, expected)
	}
}

func TestSetErrsNil(t *testing.T) {
	typed := errors.New("mock err")
	key := WrapError(typed, &specs.OnError{})

	collection := Errs{}
	collection.Set(key, nil)

	result := collection.Get(nil)
	if result != nil {
		t.Errorf("unexpected result %+v, expected %+v", result, nil)
	}
}

func TestForwardNewMeta(t *testing.T) {
	schema := specs.Header{
		"key": &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		},
	}

	forward := &Forward{
		Schema: schema,
	}

	ctx := logger.WithLogger(broker.NewBackground())
	forward.NewMeta(ctx, "mock")

	if forward.Meta == nil {
		t.Fatal("meta manager not set")
	}
}

func TestForwardNewMetaNil(t *testing.T) {
	forward := &Forward{}
	ctx := logger.WithLogger(broker.NewBackground())
	forward.NewMeta(ctx, "mock")

	if forward.Meta != nil {
		t.Fatal("unexpected meta manager")
	}
}

func TestEndpointNewCodec(t *testing.T) {
	typed := errors.New("mock error")
	endpoint := &Endpoint{
		Request: NewObject(nil, nil, nil),
		Errs:    Errs{},
		Flow: &MockFlowManager{
			errs: []Error{
				WrapError(typed, &specs.OnError{}),
			},
		},
		Response: NewObject(nil, nil, nil),
	}

	ctx := logger.WithLogger(broker.NewBackground())
	codec := json.NewConstructor()
	err := endpoint.NewCodec(ctx, codec, codec)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndpointNewCodecNil(t *testing.T) {
	endpoint := &Endpoint{}
	ctx := logger.WithLogger(broker.NewBackground())
	codec := json.NewConstructor()
	err := endpoint.NewCodec(ctx, codec, codec)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndpointNewCodecRequestErr(t *testing.T) {
	typed := errors.New("mock error")
	endpoint := &Endpoint{
		Request: NewObject(&specs.ParameterMap{}, &specs.Property{}, &specs.Property{}),
		Errs:    Errs{},
		Flow: &MockFlowManager{
			errs: []Error{
				WrapError(typed, &specs.OnError{}),
			},
		},
		Response: NewObject(nil, nil, nil),
	}

	ctx := logger.WithLogger(broker.NewBackground())
	codec := &MockCodecConstructor{err: typed}

	err := endpoint.NewCodec(ctx, codec, codec)
	if err == nil {
		t.Fatal("unexpected pass")
	}
}

func TestEndpointNewCodecResponseErr(t *testing.T) {
	typed := errors.New("mock error")
	endpoint := &Endpoint{
		Request: NewObject(nil, nil, nil),
		Errs:    Errs{},
		Flow: &MockFlowManager{
			errs: []Error{
				WrapError(typed, &specs.OnError{}),
			},
		},
		Response: NewObject(&specs.ParameterMap{}, &specs.Property{}, &specs.Property{}),
	}

	ctx := logger.WithLogger(broker.NewBackground())
	codec := &MockCodecConstructor{err: typed}

	err := endpoint.NewCodec(ctx, codec, codec)
	if err == nil {
		t.Fatal("unexpected pass")
	}
}

func TestEndpointNewCodecOnErrorErr(t *testing.T) {
	typed := errors.New("mock error")
	endpoint := &Endpoint{
		Request: NewObject(nil, nil, nil),
		Errs:    Errs{},
		Flow: &MockFlowManager{
			errs: []Error{
				WrapError(typed, &specs.OnError{
					Response: &specs.ParameterMap{},
				}),
			},
		},
		Response: NewObject(nil, nil, nil),
	}

	ctx := logger.WithLogger(broker.NewBackground())
	codec := &MockCodecConstructor{err: typed}

	err := endpoint.NewCodec(ctx, codec, codec)
	if err == nil {
		t.Fatal("unexpected pass")
	}
}
