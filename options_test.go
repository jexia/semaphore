package semaphore

import (
	"testing"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestWithFlowsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.FlowsManifest, error) { return nil, nil }

	result, err := New(WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.FlowResolvers) != 1 {
		t.Fatal("unexpected result expected flow resolver to be set")
	}
}

func TestWithMultipleFlowsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.FlowsManifest, error) { return nil, nil }

	result, err := New(WithFlows(resolver), WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.FlowResolvers) != 2 {
		t.Fatal("unexpected result expected multiple flow resolvers to be set")
	}
}

func TestWithServicesOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.ServicesManifest, error) { return nil, nil }

	result, err := New(WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.ServiceResolvers) != 1 {
		t.Fatal("unexpected result expected service resolver to be set")
	}
}

func TestWithMultipleServicesOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.ServicesManifest, error) { return nil, nil }

	result, err := New(WithServices(resolver), WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.ServiceResolvers) != 2 {
		t.Fatal("unexpected result expected multiple service resolvers to be set")
	}
}

func TestWithEndpointsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.EndpointsManifest, error) { return nil, nil }

	result, err := New(WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.EndpointResolvers) != 1 {
		t.Fatal("unexpected result expected endpoint resolver to be set")
	}
}

func TestWithMultipleEndpointsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.EndpointsManifest, error) { return nil, nil }

	result, err := New(WithEndpoints(resolver), WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.EndpointResolvers) != 2 {
		t.Fatal("unexpected result expected multiple endpoints resolvers to be set")
	}
}

func TestWithLogLevel(t *testing.T) {
	_, err := New(WithLogLevel(logger.Core, "debug"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithInvalidLogLevel(t *testing.T) {
	_, err := New(WithLogLevel(logger.Core, "unknown"))
	if err != nil {
		t.Fatal(err)
	}
}
