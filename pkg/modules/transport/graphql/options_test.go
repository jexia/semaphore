package graphql

import (
	"context"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

type MockFlow struct {
	name string
}

func (flow *MockFlow) NewStore() references.Store                          { return nil }
func (flow *MockFlow) GetName() string                                     { return flow.name }
func (flow *MockFlow) Errors() []transport.Error                           { return nil }
func (flow *MockFlow) Do(ctx context.Context, refs references.Store) error { return nil }
func (flow *MockFlow) Wait()                                               {}

func TestParseEndpointOptionsNil(t *testing.T) {
	ParseEndpointOptions(nil)
}

func TestParseEndpointOptionsNilFlow(t *testing.T) {
	ParseEndpointOptions(&transport.Endpoint{})
}

func TestParseEndpointOptionsDefault(t *testing.T) {
	name := "mock"
	base := QueryObject

	endpoint := transport.Endpoint{
		Flow: &MockFlow{name: name},
	}

	result, err := ParseEndpointOptions(&endpoint)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != name {
		t.Errorf("unexpected name %s, expected %s", result.Name, name)
	}

	if result.Path != name {
		t.Errorf("unexpected path %s, expected %s", result.Path, name)
	}

	if result.Base != base {
		t.Errorf("unexpected base %s, expected %s", result.Base, base)
	}
}

func TestParseEndpointOptions(t *testing.T) {
	name := "mock"
	base := MutationObject

	endpoint := transport.Endpoint{
		Flow: &MockFlow{name: name},
		Options: specs.Options{
			PathOption: name,
			BaseOption: base,
			NameOption: name,
		},
	}

	result, err := ParseEndpointOptions(&endpoint)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != name {
		t.Errorf("unexpected name %s, expected %s", result.Name, name)
	}

	if result.Path != name {
		t.Errorf("unexpected path %s, expected %s", result.Path, name)
	}

	if result.Base != base {
		t.Errorf("unexpected base %s, expected %s", result.Base, base)
	}
}

func TestParseEndpointOptionsBaseErr(t *testing.T) {
	name := "mock"
	base := "unknown"

	expected := "unknown base 'unknown', expected query or mutation"
	endpoint := transport.Endpoint{
		Flow: &MockFlow{name: name},
		Options: specs.Options{
			BaseOption: base,
		},
	}

	_, err := ParseEndpointOptions(&endpoint)
	if err.Error() != expected {
		t.Fatalf("unexpected error %+v, expected %s", err, expected)
	}
}
