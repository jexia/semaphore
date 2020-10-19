package metadata

import "testing"

func TestWithValue(t *testing.T) {
	type typed int32

	key := typed(1)
	value := "hello world"

	meta := WithValue(nil, key, value)
	if meta == nil {
		t.Fatal("unexpected nil meta")
	}

	result := meta.Value(key)
	if result == nil {
		t.Fatal("unexpected empty result")
	}

	if result != value {
		t.Fatalf("unexpected result %+v, expected %+v", result, value)
	}
}

func TestMultipleWithValue(t *testing.T) {
	type typed int32

	parent := typed(1)
	child := typed(2)

	expected := "hello world"
	unexpected := ""

	meta := WithValue(WithValue(nil, child, expected), parent, unexpected)
	if meta == nil {
		t.Fatal("unexpected nil meta")
	}

	result := meta.Value(child)
	if result == nil {
		t.Fatal("unexpected empty result")
	}

	if result != expected {
		t.Fatalf("unexpected result %+v, expected %+v", result, expected)
	}
}

func TestNilValue(t *testing.T) {
	type typed int32
	key := typed(1)

	var meta *Meta
	result := meta.Value(key)
	if result != nil {
		t.Fatalf("unexpected result %+v, expected nil value", result)
	}
}
