package endpoints

import (
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestNewOptionsNil(t *testing.T) {
	NewOptions(nil)
}

func TestNewOptionsWithFunctions(t *testing.T) {
	expected := functions.Collection{
		nil: nil,
	}

	options := NewOptions(WithFunctions(expected))
	if len(options.stack) != len(expected) {
		t.Fatalf("unexpected result %+v, expected %+v", options.stack, expected)
	}
}

func TestNewOptionsWithServices(t *testing.T) {
	expected := specs.ServiceList{nil}

	options := NewOptions(WithServices(expected))
	if len(options.services) != len(expected) {
		t.Fatalf("unexpected result %+v, expected %+v", options.services, expected)
	}
}

func TestNewOptionsWithCore(t *testing.T) {
	expected := semaphore.Options{
		Functions: functions.Custom{
			"": nil,
		},
	}

	options := NewOptions(WithCore(expected))
	if len(options.Functions) != len(expected.Functions) {
		t.Fatalf("unexpected result %+v, expected %+v", options.Options, expected)
	}
}

func TestNewTransportersNil(t *testing.T) {
	Transporters(nil, nil, nil, nil)
}

func TestNewTransportersEmpty(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	endpoints := specs.EndpointList{}
	flows := specs.FlowListInterface{}

	list, err := Transporters(ctx, endpoints, flows)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("unexpected empty list")
	}
}

func TestNewTransporters(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	endpoints := specs.EndpointList{
		&specs.Endpoint{
			Flow:     "mock",
			Listener: "http",
		},
	}

	flows := specs.FlowListInterface{
		&specs.Flow{
			Name: "mock",
		},
	}

	list, err := Transporters(ctx, endpoints, flows)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("unexpected empty list")
	}
}

func TestNewTransportersErr(t *testing.T) {
	expected := "failed to construct flow: nil flow manager"

	ctx := logger.WithLogger(broker.NewBackground())
	endpoints := specs.EndpointList{
		&specs.Endpoint{
			Flow:     "mock",
			Listener: "http",
		},
	}

	_, err := Transporters(ctx, endpoints, nil)
	if err.Error() != expected {
		t.Fatalf("unexpected err %s, expected %s", err.Error(), expected)
	}
}

func TestNewTransportersProxy(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	endpoints := specs.EndpointList{
		&specs.Endpoint{
			Flow:     "mock",
			Listener: "http",
		},
	}

	flows := specs.FlowListInterface{
		&specs.Proxy{
			Name: "mock",
			Forward: &specs.Call{
				Service: "mock",
				Request: &specs.ParameterMap{},
			},
		},
	}

	services := specs.ServiceList{
		&specs.Service{
			FullyQualifiedName: "mock",
		},
	}

	list, err := Transporters(ctx, endpoints, flows, WithServices(services))
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("unexpected empty list")
	}
}

func TestNewTransportersProxyUnkownService(t *testing.T) {
	expected := "failed to construct flow caller: unknown service 'mock'"

	ctx := logger.WithLogger(broker.NewBackground())
	endpoints := specs.EndpointList{
		&specs.Endpoint{
			Flow:     "mock",
			Listener: "http",
		},
	}

	flows := specs.FlowListInterface{
		&specs.Proxy{
			Name: "mock",
			Forward: &specs.Call{
				Service: "mock",
			},
		},
	}

	_, err := Transporters(ctx, endpoints, flows)
	if err.Error() != expected {
		t.Fatalf("unexpected err %s, expected %s", err.Error(), expected)
	}
}
