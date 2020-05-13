package maestro

import (
	"testing"

	"github.com/jexia/maestro/pkg/constructor"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
)

func TestAfterConstructorOption(t *testing.T) {
	fn := func(next constructor.AfterConstructor) constructor.AfterConstructor { return next }

	result, err := New(AfterConstructor(fn))
	if err != nil {
		t.Fatal(err)
	}

	if result.Options.AfterConstructor == nil {
		t.Fatal("unexpected result expected after constructor to be set")
	}
}

func TestMultipleAfterConstructorOption(t *testing.T) {
	fn := func(next constructor.AfterConstructor) constructor.AfterConstructor { return next }

	result, err := New(AfterConstructor(fn), AfterConstructor(fn))
	if err != nil {
		t.Fatal(err)
	}

	if result.Options.AfterConstructor == nil {
		t.Fatal("unexpected result expected after constructor to be set")
	}
}

func TestWithFlowsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.FlowsManifest, error) { return nil, nil }

	result, err := New(WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.Flows) != 1 {
		t.Fatal("unexpected result expected flow resolver to be set")
	}
}

func TestWithMultipleFlowsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.FlowsManifest, error) { return nil, nil }

	result, err := New(WithFlows(resolver), WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.Flows) != 2 {
		t.Fatal("unexpected result expected multiple flow resolvers to be set")
	}
}

func TestWithServicesOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.ServicesManifest, error) { return nil, nil }

	result, err := New(WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.Services) != 1 {
		t.Fatal("unexpected result expected service resolver to be set")
	}
}

func TestWithMultipleServicesOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.ServicesManifest, error) { return nil, nil }

	result, err := New(WithServices(resolver), WithServices(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.Services) != 2 {
		t.Fatal("unexpected result expected multiple service resolvers to be set")
	}
}

func TestWithEndpointsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.EndpointsManifest, error) { return nil, nil }

	result, err := New(WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.Endpoints) != 1 {
		t.Fatal("unexpected result expected endpoint resolver to be set")
	}
}

func TestWithMultipleEndpointsOption(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.EndpointsManifest, error) { return nil, nil }

	result, err := New(WithEndpoints(resolver), WithEndpoints(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Options.Endpoints) != 2 {
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
	_, err := New(WithLogLevel(logger.Core, "unkown"))
	if err != nil {
		t.Fatal(err)
	}
}