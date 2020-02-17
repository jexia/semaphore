package refs

import (
	"encoding/json"
	"testing"

	"github.com/jexia/maestro/specs"
)

func BenchmarkSimpleUnmarshal(b *testing.B) {
	input := []byte(`{"message":"hello world"}`)

	data := map[string]interface{}{}
	json.Unmarshal(input, &data)

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewStore(len(data))
		store.StoreValues("input", "", data)
	}
}

func BenchmarkRepeatedUnmarshal(b *testing.B) {
	input := []byte(`{"details":[{"message":"hello world"}]}`)

	data := map[string]interface{}{}
	json.Unmarshal(input, &data)

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewStore(len(data))
		store.StoreValues("input", "", data)
	}
}

func BenchmarkNestedUnmarshal(b *testing.B) {
	input := []byte(`{"details":{"name":"john"}}`)

	data := map[string]interface{}{}
	json.Unmarshal(input, &data)

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewStore(len(data))
		store.StoreValues("input", "", data)
	}
}

func BenchmarkComplexUnmarshal(b *testing.B) {
	input := []byte(`{"message":"hello","details":{"name":"john"},"collection":[{"name":"john"}]}`)

	data := map[string]interface{}{}
	json.Unmarshal(input, &data)

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewStore(len(data))
		store.StoreValues("input", "", data)
	}
}

func TestStoreReference(t *testing.T) {
	store := NewStore(10)

	resource := "input"
	ref := Reference{
		Path:  "test",
		Value: "hello world",
	}

	store.StoreReference(resource, ref)
	result := store.Load(resource, ref.Path)
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Path != ref.Path {
		t.Fatalf("unexpected path %s, expected %s", result.Path, ref.Path)
	}

	if result.Value != ref.Value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, ref.Value)
	}
}

func TestStoreValues(t *testing.T) {
	store := NewStore(1)

	target := "message"
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		target: value,
	}

	store.StoreValues(resource, "", values)
	result := store.Load(resource, target)
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Path != target {
		t.Fatalf("unexpected path %s, expected %s", result.Path, target)
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}

func TestStoreNestedValues(t *testing.T) {
	store := NewStore(1)

	nested := "nested"
	key := "message"
	target := specs.JoinPath(nested, key)
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: map[string]interface{}{
			key: value,
		},
	}

	store.StoreValues(resource, "", values)
	result := store.Load(resource, target)
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Path != target {
		t.Fatalf("unexpected path %s, expected %s", result.Path, target)
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}

func TestStoreRepeatedValues(t *testing.T) {
	store := NewStore(1)

	nested := "nested"
	key := "message"
	target := specs.JoinPath(nested, key)
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: []interface{}{
			map[string]interface{}{
				key: value,
			},
		},
	}

	store.StoreValues(resource, "", values)
	result := store.Load(resource, nested)
	if result == nil {
		t.Fatal("did not return reference")
	}

	if len(result.Repeated) != 1 {
		t.Fatalf("unexpected repeated length %d, expected 1", len(result.Repeated))
	}

	if result.Path != nested {
		t.Fatalf("unexpected repeated reference path %s, expected %s", result.Path, nested)
	}

	repeating := result.Repeated[0]
	result = repeating.Load(resource, target)
	if result == nil {
		t.Fatal("did not return repeating reference")
	}

	if result.Path != target {
		t.Fatalf("unexpected path %s, expected %s", result.Path, target)
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}
