package providers

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestResolve(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResolveBeforeMiddleware(t *testing.T) {
	counter := 0
	options := Options{
		Options: semaphore.Options{
			BeforeConstructor: func(*broker.Context, functions.Collection, semaphore.Options) error {
				counter++
				return nil
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != nil {
		t.Fatal(err)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}

}

func TestResolveBeforeMiddlewareErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	options := Options{
		Options: semaphore.Options{
			BeforeConstructor: func(*broker.Context, functions.Collection, semaphore.Options) error {
				counter++
				return expected
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != expected {
		t.Fatalf("unexpected err %+v, expected %+v", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}

}

func TestResolveAfterMiddleware(t *testing.T) {
	counter := 0
	options := Options{
		AfterConstructor: func(*broker.Context, specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas) error {
			counter++
			return nil
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != nil {
		t.Fatal(err)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}
}

func TestResolveAfterMiddlewareErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	options := Options{
		AfterConstructor: func(*broker.Context, specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas) error {
			counter++
			return expected
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != expected {
		t.Fatalf("unexpected err %+v, expected %+v", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}
}

func TestAfterConstructorOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) AfterConstructorHandler {
		return func(next AfterConstructor) AfterConstructor {
			return func(ctx *broker.Context, flow specs.FlowListInterface, endpoints specs.EndpointList, services specs.ServiceList, schemas specs.Schemas) error {
				*i++
				return next(ctx, flow, endpoints, services, schemas)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := []Option{WithAfterConstructor(fn(&result))}

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := []Option{WithAfterConstructor(fn(&result)), WithAfterConstructor(fn(&result))}

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, semaphore.Options{}, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.AfterConstructor == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			err = client.AfterConstructor(nil, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestFlowResolverErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	options := Options{
		Options: semaphore.Options{
			FlowResolvers: providers.FlowsResolvers{
				func(*broker.Context) (specs.FlowListInterface, error) {
					counter++
					return nil, expected
				},
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != expected {
		t.Fatalf("unexpected err %+v, expected %+v", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}
}

func TestEndpointResolverErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	options := Options{
		EndpointResolvers: providers.EndpointResolvers{
			func(*broker.Context) (specs.EndpointList, error) {
				counter++
				return nil, expected
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != expected {
		t.Fatalf("unexpected err %+v, expected %+v", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}
}

func TestSchemaResolverErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	options := Options{
		SchemaResolvers: providers.SchemaResolvers{
			func(*broker.Context) (specs.Schemas, error) {
				counter++
				return nil, expected
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != expected {
		t.Fatalf("unexpected err %+v, expected %+v", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}
}

func TestServicesResolverErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	options := Options{
		ServiceResolvers: providers.ServiceResolvers{
			func(*broker.Context) (specs.ServiceList, error) {
				counter++
				return nil, expected
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := Resolve(ctx, functions.Collection{}, options)
	if err != expected {
		t.Fatalf("unexpected err %+v, expected %+v", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected %d", counter, 1)
	}
}
