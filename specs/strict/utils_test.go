package strict

import (
	"testing"

	"github.com/jexia/maestro/specs"
)

func TestGetService(t *testing.T) {
	tests := map[string]string{
		"service.Call": "service",
		"service":      "service",
		"service.":     "service",
	}

	for input, expected := range tests {
		result := GetService(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestGetMethod(t *testing.T) {
	tests := map[string]string{
		"service.Call":     "Call",
		"service":          "",
		"service.Call.sub": "Call.sub",
	}

	for input, expected := range tests {
		result := GetMethod(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestGetSchemaService(t *testing.T) {
	manifest := &specs.Manifest{
		Services: []*specs.Service{
			{
				Alias:  "service",
				Schema: "servicing",
			},
			{
				Alias:  "call",
				Schema: "calling",
			},
		},
	}

	tests := map[string]string{
		"service": manifest.Services[0].Schema,
		"call":    manifest.Services[1].Schema,
		"":        "",
	}

	for input, expected := range tests {
		result := GetSchemaService(manifest, input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}
