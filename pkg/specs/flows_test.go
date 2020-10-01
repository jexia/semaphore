package specs

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs/metadata"
)

func TestFlowListInterfaceAppend(t *testing.T) {
	flows := FlowListInterface{}
	flows.Append(FlowListInterface{&Flow{}, &Flow{}})

	if len(flows) != 2 {
		t.Fatalf("unexpected length %+v, expected 2", len(flows))
	}
}

func TestFlowListInterfaceGet(t *testing.T) {
	expected := &Flow{Name: "expected"}

	flows := FlowListInterface{&Flow{Name: "unexpected"}, expected}
	result := flows.Get("expected")

	if result != expected {
		t.Fatalf("unexpected result %+v, expected %+v", result, expected)
	}
}

func TestFlowListInterfaceGetUnknown(t *testing.T) {
	flows := FlowListInterface{}
	result := flows.Get("expected")

	if result != nil {
		t.Fatalf("unexpected result %+v, expected nil value", result)
	}
}

func TestFlowListGet(t *testing.T) {
	expected := &Flow{Name: "expected"}

	flows := FlowList{&Flow{Name: "unexpected"}, expected}
	result := flows.Get("expected")

	if result != expected {
		t.Fatalf("unexpected result %+v, expected %+v", result, expected)
	}
}

func TestFlowListGetUnknown(t *testing.T) {
	flows := FlowList{}
	result := flows.Get("expected")

	if result != nil {
		t.Fatalf("unexpected result %+v, expected nil value", result)
	}
}

func TestFlowInterface(t *testing.T) {
	flow := &Flow{
		Meta:    metadata.WithValue(nil, nil, nil),
		Name:    "sample",
		Input:   &ParameterMap{},
		Nodes:   NodeList{&Node{ID: "first"}},
		Output:  &ParameterMap{},
		OnError: &OnError{},
	}

	if flow.GetName() != "sample" {
		t.Error("unexpected result, expected name to be set")
	}

	if flow.GetInput() == nil {
		t.Error("unexpected result, expected input to be set")
	}

	if flow.GetNodes() == nil {
		t.Error("unexpected result, expected nodes to be set")
	}

	if flow.GetNodes().Get("first") == nil {
		t.Error("unexpected result, expected first node to be set")
	}

	flow.SetNodes(NodeList{&Node{ID: "second"}})

	if flow.GetNodes() == nil {
		t.Error("unexpected result, expected nodes to be set")
	}

	if flow.GetNodes().Get("second") == nil {
		t.Error("unexpected result, expected second node to be set")
	}

	if flow.GetOutput() == nil {
		t.Error("unexpected result, expected output to be set")
	}

	if flow.GetOnError() == nil {
		t.Error("unexpected result, expected on error to be set")
	}

	if flow.GetMeta() == nil {
		t.Error("unexpected result, expected meta to be set")
	}

	if flow.GetForward() != nil {
		t.Error("unexpected result, expected forward to be a nil value")
	}
}

func TestProxyListGet(t *testing.T) {
	expected := &Proxy{Name: "expected"}

	flows := ProxyList{&Proxy{Name: "unexpected"}, expected}
	result := flows.Get("expected")

	if result != expected {
		t.Fatalf("unexpected result %+v, expected %+v", result, expected)
	}
}

func TestProxyListGetUnknown(t *testing.T) {
	flows := ProxyList{}
	result := flows.Get("expected")

	if result != nil {
		t.Fatalf("unexpected result %+v, expected nil value", result)
	}
}

func TestProxyInterface(t *testing.T) {
	flow := &Proxy{
		Meta:    metadata.WithValue(nil, nil, nil),
		Name:    "sample",
		Input:   &ParameterMap{},
		Nodes:   NodeList{&Node{ID: "first"}},
		Forward: &Call{},
		OnError: &OnError{},
	}

	if flow.GetName() != "sample" {
		t.Error("unexpected result, expected name to be set")
	}

	if flow.GetInput() == nil {
		t.Error("unexpected result, expected input to be set")
	}

	if flow.GetNodes().Get("first") == nil {
		t.Error("unexpected result, expected first node to be set")
	}

	flow.SetNodes(NodeList{&Node{ID: "second"}})

	if flow.GetNodes() == nil {
		t.Error("unexpected result, expected nodes to be set")
	}

	if flow.GetNodes().Get("second") == nil {
		t.Error("unexpected result, expected second node to be set")
	}

	if flow.GetOutput() != nil {
		t.Error("unexpected result, expected output to be a nil value")
	}

	if flow.GetOnError() == nil {
		t.Error("unexpected result, expected on error to be set")
	}

	if flow.GetMeta() == nil {
		t.Error("unexpected result, expected meta to be set")
	}

	if flow.GetForward() == nil {
		t.Error("unexpected result, expected forward to be set")
	}
}

func TestNode(t *testing.T) {
	node := &Node{
		Meta:    metadata.WithValue(nil, nil, nil),
		Name:    "sample",
		OnError: &OnError{},
	}

	if node.GetOnError() == nil {
		t.Error("unexpected result, expected on error to be set")
	}
}

func TestNodeListGet(t *testing.T) {
	nodes := &NodeList{&Node{ID: "first"}, &Node{ID: "second"}}
	result := nodes.Get("second")
	if result == nil {
		t.Error("unexpected empty result")
	}
}

func TestNodeListGetUnknown(t *testing.T) {
	nodes := &NodeList{&Node{ID: "first"}, &Node{ID: "second"}}
	result := nodes.Get("unknown")
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}

func TestOnErrorClone(t *testing.T) {
	err := &OnError{
		Meta:     metadata.WithValue(nil, nil, nil),
		Response: &ParameterMap{},
		Status:   &Property{},
		Message:  &Property{},
		Params: map[string]*Property{
			"sample": {},
		},
	}

	result := err.Clone()
	if result == nil {
		t.Error("unexpected result, expected on error clone to be returned")
	}

	if result.Meta != err.Meta {
		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, err.Meta)
	}

	if result.Response == nil || result.Response == err.Response {
		t.Errorf("unexpected response %+v", result.Response)
	}

	if result.Status == nil || result.Status == err.Status {
		t.Errorf("unexpected status %+v", result.Status)
	}

	if result.Message == nil || result.Message == err.Message {
		t.Errorf("unexpected message %+v", result.Message)
	}

	if result.Message == nil || result.Message == err.Message {
		t.Errorf("unexpected message %+v", result.Message)
	}

	if result.Params == nil || len(result.Params) != len(err.Params) {
		t.Errorf("unexpected params %+v", result.Params)
	}
}

func TestOnErrorCloneNilValue(t *testing.T) {
	var err *OnError
	result := err.Clone()
	if result != nil {
		t.Fatalf("unexpected result %+v, expected nil value", result)
	}
}

func TestOnErrorInterface(t *testing.T) {
	err := &OnError{
		Meta:     metadata.WithValue(nil, nil, nil),
		Response: &ParameterMap{},
		Status:   &Property{},
		Message:  &Property{},
		Params: map[string]*Property{
			"sample": {},
		},
	}

	if err.GetResponse() == nil {
		t.Error("unexpected result, expected response to be returned")
	}

	if err.GetStatusCode() == nil {
		t.Error("unexpected result, expected status code to be returned")
	}

	if err.GetMessage() == nil {
		t.Error("unexpected result, expected message to be returned")
	}
}

func TestOnErrorInterfaceNil(t *testing.T) {
	var err *OnError

	if err.GetResponse() != nil {
		t.Errorf("unexpected result %+v, expected response nil value", err.GetResponse())
	}

	if err.GetStatusCode() != nil {
		t.Errorf("unexpected result %+v, expected status code nil value", err.GetStatusCode())
	}

	if err.GetMessage() != nil {
		t.Errorf("unexpected result %+v, expected message nil value", err.GetMessage())
	}
}

func TestDependenciesAppend(t *testing.T) {
	deps := Dependencies{}
	deps = deps.Append(Dependencies{"input": nil})

	if len(deps) != 1 {
		t.Fatalf("unexpected dependencies length %d, expected 1", len(deps))
	}

	deps = deps.Append(Dependencies{"first": nil, "second": nil})

	if len(deps) != 3 {
		t.Fatalf("unexpected dependencies length %d, expected 3", len(deps))
	}

	deps = deps.Append(Dependencies{"first": nil})

	if len(deps) != 3 {
		t.Fatalf("unexpected dependencies length %d, dependency length should have stayed at 3", len(deps))
	}
}
