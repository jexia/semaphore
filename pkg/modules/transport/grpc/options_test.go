package grpc

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
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
		{
			ServiceOption: "name",
		}: true,
		{
			MethodOption: "name",
		}: true,
		{
			PackageOption: "name",
		}: true,
		{
			ServiceOption: "name",
			MethodOption:  "name",
			PackageOption: "name",
		}: true,
		{
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
	tests := map[*specs.Options]bool{}

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
