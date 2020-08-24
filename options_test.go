package semaphore

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport/http"
)

func TestWithFlowsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, nil }

	result, err := NewOptions(ctx, WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.FlowResolvers) != 1 {
		t.Fatal("unexpected result expected flow resolver to be set")
	}
}

func TestWithMultipleFlowsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, nil }

	result, err := NewOptions(ctx, WithFlows(resolver), WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.FlowResolvers) != 2 {
		t.Fatal("unexpected result expected multiple flow resolvers to be set")
	}
}

func TestWithServicesOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.ServiceList, error) { return nil, nil }

	result, err := NewOptions(ctx, WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.ServiceResolvers) != 1 {
		t.Fatal("unexpected result expected service resolver to be set")
	}
}

func TestWithMultipleServicesOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.ServiceList, error) { return nil, nil }

	result, err := NewOptions(ctx, WithServices(resolver), WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.ServiceResolvers) != 2 {
		t.Fatal("unexpected result expected multiple service resolvers to be set")
	}
}

func TestWithEndpointsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.EndpointList, error) { return nil, nil }

	result, err := NewOptions(ctx, WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.EndpointResolvers) != 1 {
		t.Fatal("unexpected result expected endpoint resolver to be set")
	}
}

func TestWithMultipleEndpointsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.EndpointList, error) { return nil, nil }

	result, err := NewOptions(ctx, WithEndpoints(resolver), WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.EndpointResolvers) != 2 {
		t.Fatal("unexpected result expected multiple endpoints resolvers to be set")
	}
}

func TestWithLogLevel(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithLogLevel("*", "debug"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithFunctions(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithFunctions(functions.Custom{"mock": nil}))
	if err != nil {
		t.Fatal(err)
	}

	if options.Functions == nil {
		t.Fatal("functions not set")
	}

	_, has := options.Functions["mock"]
	if !has {
		t.Fatal("mock function does not exist")
	}
}

func TestWithCaller(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithCaller(http.NewCaller()))
	if err != nil {
		t.Fatal(err)
	}

	if options.Callers == nil {
		t.Fatal("callers not set")
	}

	caller := options.Callers.Get("http")
	if caller == nil {
		t.Fatal("HTTP caller does not exist")
	}
}

func TestWithListener(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithListener(http.NewListener(":0")))
	if err != nil {
		t.Fatal(err)
	}

	if options.Listeners == nil {
		t.Fatal("listeners not set")
	}

	listener := options.Listeners.Get("http")
	if listener == nil {
		t.Fatal("HTTP listener does not exist")
	}
}

func TestWithCodec(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithCodec(json.NewConstructor()))
	if err != nil {
		t.Fatal(err)
	}

	if options.Codec == nil {
		t.Fatal("codecs not set")
	}

	codec := options.Codec.Get("json")
	if codec == nil {
		t.Fatal("JSON codec does not exist")
	}
}

func TestWithSchema(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithSchema(nil))
	if err != nil {
		t.Fatal(err)
	}

	if options.SchemaResolvers == nil {
		t.Fatal("schema resolves not set")
	}

	if len(options.SchemaResolvers) != 1 {
		t.Fatal("schema resolver not set")
	}
}

func TestWithInvalidLogLevel(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithLogLevel("*", "unknown"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptions(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsNil(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsMiddleware(t *testing.T) {
	middleware := func(*broker.Context) ([]Option, error) {
		options := []Option{
			WithLogLevel("*", "warn"),
		}

		return options, nil
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithMiddleware(middleware))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsMiddlewareErr(t *testing.T) {
	expected := errors.New("unexpected err")
	middleware := func(*broker.Context) ([]Option, error) {
		return nil, expected
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithMiddleware(middleware))
	if err != expected {
		t.Fatalf("unexpected err (%+v), expected (%+v)", err, expected)
	}
}
