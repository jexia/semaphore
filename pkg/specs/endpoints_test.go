package specs

import "testing"

func TestEndpointsAppend(t *testing.T) {
	endpoints := EndpointList{}
	endpoints.Append(EndpointList{&Endpoint{}, &Endpoint{}})

	if len(endpoints) != 2 {
		t.Fatalf("unexpected length %+v, expected 2", len(endpoints))
	}
}

func TestEndpointsGet(t *testing.T) {
	var result []*Endpoint
	endpoints := EndpointList{&Endpoint{Flow: "first"}, &Endpoint{Flow: "first"}, &Endpoint{Flow: "second"}}

	result = endpoints.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result")
	}

	if len(result) != 1 {
		t.Fatalf("unexpected result length returned %d, expected 1", len(result))
	}

	result = endpoints.Get("first")
	if result == nil {
		t.Fatal("unexpected empty result")
	}

	if len(result) != 2 {
		t.Fatalf("unexpected result length returned %d, expected 1", len(result))
	}
}
