package hcl

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
)

func TestResolveService(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/services.pass.hcl",
		"pass":  "./tests/*.pass.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDiscoveryClientsResolver(t *testing.T) {
	var (
		resolver = DiscoveryClientsResolver("./tests/discovery.pass.hcl")
		ctx      = logger.WithLogger(broker.NewBackground())
	)

	clients, err := resolver(ctx)
	if err != nil {
		t.Fatalf("DiscoveryClientsResolver()() error: %v", err)
	}

	t.Run("build client with provider taken from the block name", func(t *testing.T) {
		consul, ok := clients["consul"]
		if !ok {
			t.Errorf("expect clients to include 'consul'")
			return
		}

		if consul.Provider() != "consul" {
			t.Errorf("expect configuration 'consul' to have provider 'consul', but it is '%s'", consul.Provider())
			return
		}
	})

	t.Run("build client with custom name and provider taken from the block option", func(t *testing.T) {
		consul, ok := clients["foobar"]
		if !ok {
			t.Errorf("expect clients to include 'foobar'")
			return
		}

		if consul.Provider() != "consul" {
			t.Errorf("expect configuration 'foobar' to have provider 'consul', but it is '%s'", consul.Provider())
			return
		}
	})
}

func TestResolveServiceFail(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/services.fail.hcl",
		"pass":  "./tests/*.fail.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestResolveIncludeFail(t *testing.T) {
	tests := map[string]string{
		"include": "./tests/include.fail.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestPathError(t *testing.T) {
	type fields struct {
		Path string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Path: "/semaphore"},
			"unable to resolve path, no files found '/semaphore'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrPathNotFound{
				Path: "/semaphore",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveInclude(t *testing.T) {
	tests := map[string]string{
		"include": "./tests/include.pass.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestResolveIncludeNoFiles(t *testing.T) {
	tests := map[string]string{
		"include": "./tests/unknown.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err == nil {
				t.Fatal("unexpected pass expected resolver to return a error")
			}
		})
	}
}

func TestResolveFlows(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/flows.pass.hcl",
		"pass":  "./tests/*.pass.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := FlowsResolver(path)

			_, err := resolver(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestResolveFlowsFail(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/flows.fail.hcl",
		"pass":  "./tests/*.fail.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := FlowsResolver(path)

			_, err := resolver(ctx)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestResolveEndpoints(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/endpoints.pass.hcl",
		"pass":  "./tests/*.pass.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := EndpointsResolver(path)

			_, err := resolver(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestResolveEndpointsFail(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/endpoints.fail.hcl",
		"pass":  "./tests/*.fail.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			resolver := EndpointsResolver(path)

			_, err := resolver(ctx)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestResolveOptions(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/options.pass.hcl",
		"pass":  "./tests/*.pass.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := GetOptions(ctx, path)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
