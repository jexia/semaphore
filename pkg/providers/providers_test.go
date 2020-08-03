package providers

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestFlowResolvers(t *testing.T) {
	resolver := func(instance.Context) (specs.FlowListInterface, error) {
		return specs.FlowListInterface{&specs.Flow{}}, nil
	}

	ctx := instance.NewContext()

	resolvers := FlowsResolvers{resolver}
	flows, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(flows) != 1 {
		t.Fatalf("unexpected flows %+v, expected 1 flow", flows)
	}
}

func TestFlowResolversErr(t *testing.T) {
	expected := errors.New("mock")
	resolver := func(instance.Context) (specs.FlowListInterface, error) {
		return nil, expected
	}

	ctx := instance.NewContext()

	resolvers := FlowsResolvers{resolver}
	_, err := resolvers.Resolve(ctx)
	if err == nil {
		t.Fatal("unexpected pass")
	}

	if err != expected {
		t.Fatalf("unexpected error %s, expected %s", err, expected)
	}
}

func TestNilFlowResolvers(t *testing.T) {
	ctx := instance.NewContext()
	resolvers := FlowsResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceResolvers(t *testing.T) {
	resolver := func(instance.Context) (specs.ServiceList, error) {
		return specs.ServiceList{&specs.Service{}}, nil
	}

	ctx := instance.NewContext()

	resolvers := ServiceResolvers{resolver}
	services, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(services) != 1 {
		t.Fatalf("unexpected services %+v, expected 1 service", services)
	}
}

func TestServiceResolversErr(t *testing.T) {
	expected := errors.New("mock")
	resolver := func(instance.Context) (specs.ServiceList, error) {
		return nil, expected
	}

	ctx := instance.NewContext()

	resolvers := ServiceResolvers{resolver}
	_, err := resolvers.Resolve(ctx)
	if err == nil {
		t.Fatal("unexpected pass")
	}

	if err != expected {
		t.Fatalf("unexpected error %s, expected %s", err, expected)
	}
}

func TestNilServiceResolvers(t *testing.T) {
	ctx := instance.NewContext()
	resolvers := ServiceResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchemaResolvers(t *testing.T) {
	resolver := func(instance.Context) (specs.Objects, error) {
		return specs.Objects{"mock": &specs.Property{}}, nil
	}

	ctx := instance.NewContext()

	resolvers := SchemaResolvers{resolver}
	schemas, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(schemas) != 1 {
		t.Fatalf("unexpected schemas %+v, expected 1 schema", schemas)
	}
}

func TestSchemaResolversErr(t *testing.T) {
	expected := errors.New("mock")
	resolver := func(instance.Context) (specs.Objects, error) {
		return nil, expected
	}

	ctx := instance.NewContext()

	resolvers := SchemaResolvers{resolver}
	_, err := resolvers.Resolve(ctx)
	if err == nil {
		t.Fatal("unexpected pass")
	}

	if err != expected {
		t.Fatalf("unexpected error %s, expected %s", err, expected)
	}
}

func TestNilSchemaResolvers(t *testing.T) {
	ctx := instance.NewContext()
	resolvers := SchemaResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndpointResolvers(t *testing.T) {
	resolver := func(instance.Context) (specs.EndpointList, error) {
		return specs.EndpointList{&specs.Endpoint{}}, nil
	}

	ctx := instance.NewContext()

	resolvers := EndpointResolvers{resolver}
	endpoints, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("unexpected endpoints %+v, expected 1 endpoint", endpoints)
	}
}

func TestEndpointResolversErr(t *testing.T) {
	expected := errors.New("mock")
	resolver := func(instance.Context) (specs.EndpointList, error) {
		return nil, expected
	}

	ctx := instance.NewContext()

	resolvers := EndpointResolvers{resolver}
	_, err := resolvers.Resolve(ctx)
	if err == nil {
		t.Fatal("unexpected pass")
	}

	if err != expected {
		t.Fatalf("unexpected error %s, expected %s", err, expected)
	}
}

func TestNilEndpointResolvers(t *testing.T) {
	ctx := instance.NewContext()
	resolvers := EndpointResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
