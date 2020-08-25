package providers

import (
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport/http"
)

func TestWithListener(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, semaphore.Options{}, WithListener(http.NewListener(":0")))
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

func TestWithSchema(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, semaphore.Options{}, WithSchema(nil))
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

func TestWithCoreOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	core := semaphore.Options{
		Functions: functions.Custom{"make": nil},
	}

	result, err := NewOptions(ctx, core)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Functions) != 1 {
		t.Fatal("unexpected result expected functions to be set")
	}
}

func TestWithServicesOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.ServiceList, error) { return nil, nil }

	result, err := NewOptions(ctx, semaphore.Options{}, WithServices(resolver))
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

	result, err := NewOptions(ctx, semaphore.Options{}, WithServices(resolver), WithServices(resolver))
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

	result, err := NewOptions(ctx, semaphore.Options{}, WithEndpoints(resolver))
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

	result, err := NewOptions(ctx, semaphore.Options{}, WithEndpoints(resolver), WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.EndpointResolvers) != 2 {
		t.Fatal("unexpected result expected multiple endpoints resolvers to be set")
	}
}

func TestNewOptions(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, semaphore.Options{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsNil(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, semaphore.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
