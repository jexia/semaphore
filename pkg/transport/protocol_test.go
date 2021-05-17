package transport

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/discovery"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/codec"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

type MockCaller struct{}

func (mock *MockCaller) Name() string { return "mock" }
func (mock *MockCaller) Dial(service *specs.Service, functions functions.Custom, options specs.Options, resolver discovery.Resolver) (Call, error) {
	return nil, nil
}

func TestGetCaller(t *testing.T) {
	list := Callers{&MockCaller{}}
	result := list.Get("mock")
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestGetCallerUnkown(t *testing.T) {
	list := Callers{&MockCaller{}}
	result := list.Get("unknown")
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}

type MockListener struct{}

func (mock *MockListener) Name() string { return "mock" }
func (mock *MockListener) Serve() error { return nil }
func (mock *MockListener) Close() error { return nil }
func (mock *MockListener) Handle(*broker.Context, []*Endpoint, map[string]codec.Constructor) error {
	return nil
}

func TestGetListener(t *testing.T) {
	list := ListenerList{&MockListener{}}
	result := list.Get("mock")
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestGetListenerUnkown(t *testing.T) {
	list := ListenerList{&MockListener{}}
	result := list.Get("unknown")
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}
