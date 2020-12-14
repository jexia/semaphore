package providers

import (
	"errors"
	"github.com/jexia/semaphore/pkg/discovery"
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestFlowResolvers(t *testing.T) {
	resolver := func(*broker.Context) (specs.FlowListInterface, error) {
		return specs.FlowListInterface{&specs.Flow{}}, nil
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	resolver := func(*broker.Context) (specs.FlowListInterface, error) {
		return nil, expected
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	ctx := logger.WithLogger(broker.NewBackground())
	resolvers := FlowsResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceResolvers(t *testing.T) {
	resolver := func(*broker.Context) (specs.ServiceList, error) {
		return specs.ServiceList{&specs.Service{}}, nil
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	resolver := func(*broker.Context) (specs.ServiceList, error) {
		return nil, expected
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	ctx := logger.WithLogger(broker.NewBackground())
	resolvers := ServiceResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchemaResolvers(t *testing.T) {
	resolver := func(*broker.Context) (specs.Schemas, error) {
		return specs.Schemas{"mock": &specs.Property{}}, nil
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	resolver := func(*broker.Context) (specs.Schemas, error) {
		return nil, expected
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	ctx := logger.WithLogger(broker.NewBackground())
	resolvers := SchemaResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndpointResolvers(t *testing.T) {
	resolver := func(*broker.Context) (specs.EndpointList, error) {
		return specs.EndpointList{&specs.Endpoint{}}, nil
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	resolver := func(*broker.Context) (specs.EndpointList, error) {
		return nil, expected
	}

	ctx := logger.WithLogger(broker.NewBackground())

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
	ctx := logger.WithLogger(broker.NewBackground())
	resolvers := EndpointResolvers{nil}
	_, err := resolvers.Resolve(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultServiceResolverClient_Resolver(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    discovery.Resolver
		wantErr bool
	}{
		{
			"",
			args{"http://localhost:3000"},
			discovery.NewPlainResolver("http://localhost:3000"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dnsServiceResolver{}
			got, err := d.Resolver(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver() got = %v, want %v", got, tt.want)
			}
		})
	}
}
