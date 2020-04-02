package grpc

import (
	"testing"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
)

func TestListenerOptions(t *testing.T) {
	tests := map[*specs.Options]bool{}

	for test, pass := range tests {
		_, err := ParseListenerOptions(*test)
		if err == nil && !pass {
			t.Fatalf("unexpected pass: %+v", test)
		}

		if err != nil && pass {
			t.Fatalf("unexpected fail: %+v", test)
		}
	}
}

func TestEndpointOptions(t *testing.T) {
	tests := map[*specs.Options]bool{
		&specs.Options{
			ServiceOption: "name",
		}: true,
		&specs.Options{
			MethodOption: "name",
		}: true,
		&specs.Options{
			PackageOption: "name",
		}: true,
		&specs.Options{
			ServiceOption: "name",
			MethodOption:  "name",
			PackageOption: "name",
		}: true,
		&specs.Options{
			"unknown": "name",
		}: true,
	}

	for test, pass := range tests {
		_, err := ParseListenerOptions(*test)
		if err == nil && !pass {
			t.Fatalf("unexpected pass: %+v", test)
		}

		if err != nil && pass {
			t.Fatalf("unexpected fail: %+v", test)
		}
	}
}

func TestCallerOptions(t *testing.T) {
	tests := map[*schema.Options]bool{}

	for test, pass := range tests {
		_, err := ParseCallerOptions(*test)
		if err == nil && !pass {
			t.Fatalf("unexpected pass: %+v", test)
		}

		if err != nil && pass {
			t.Fatalf("unexpected fail: %+v", test)
		}
	}
}
