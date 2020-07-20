package hcl

import (
	"testing"

	"github.com/jexia/semaphore/pkg/core/instance"
)

func TestResolveService(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/services.pass.hcl",
		"pass":  "./tests/*.pass.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestResolveServiceFail(t *testing.T) {
	tests := map[string]string{
		"basic": "./tests/services.fail.hcl",
		"pass":  "./tests/*.fail.hcl",
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
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
			ctx := instance.NewContext()
			resolver := ServicesResolver(path)

			_, err := resolver(ctx)
			if err == nil {
				t.Fatal("unexpected pass")
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
			ctx := instance.NewContext()
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
			ctx := instance.NewContext()
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
			ctx := instance.NewContext()
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
			ctx := instance.NewContext()
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
			ctx := instance.NewContext()
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
			ctx := instance.NewContext()
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
			_, err := GetOptions(instance.NewContext(), path)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
