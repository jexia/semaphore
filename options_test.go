package semaphore

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestWithFlowsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, nil }

	result, err := New(ctx, WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.FlowResolvers) != 1 {
		t.Fatal("unexpected result expected flow resolver to be set")
	}
}

func TestWithMultipleFlowsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, nil }

	result, err := New(ctx, WithFlows(resolver), WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.FlowResolvers) != 2 {
		t.Fatal("unexpected result expected multiple flow resolvers to be set")
	}
}

func TestWithServicesOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.ServiceList, error) { return nil, nil }

	result, err := New(ctx, WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.ServiceResolvers) != 1 {
		t.Fatal("unexpected result expected service resolver to be set")
	}
}

func TestWithMultipleServicesOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.ServiceList, error) { return nil, nil }

	result, err := New(ctx, WithServices(resolver), WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.ServiceResolvers) != 2 {
		t.Fatal("unexpected result expected multiple service resolvers to be set")
	}
}

func TestWithEndpointsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.EndpointList, error) { return nil, nil }

	result, err := New(ctx, WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.EndpointResolvers) != 1 {
		t.Fatal("unexpected result expected endpoint resolver to be set")
	}
}

func TestWithMultipleEndpointsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.EndpointList, error) { return nil, nil }

	result, err := New(ctx, WithEndpoints(resolver), WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.EndpointResolvers) != 2 {
		t.Fatal("unexpected result expected multiple endpoints resolvers to be set")
	}
}

func TestWithLogLevel(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := New(ctx, WithLogLevel("*", "debug"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithInvalidLogLevel(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := New(ctx, WithLogLevel("*", "unknown"))
	if err != nil {
		t.Fatal(err)
	}
}
