package strict

import (
	"testing"
)

func TestGetService(t *testing.T) {
	tests := map[string]string{
		"service.Call":    "service",
		"service":         "service",
		"service.":        "service",
		"schema.service.": "schema",
		"schema.service":  "schema",
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
		"service.Call":        "Call",
		"service":             "",
		"schema.service.Call": "Call",
	}

	for input, expected := range tests {
		result := GetMethod(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}
