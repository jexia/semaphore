package refs

import (
	"encoding/json"
	"testing"

	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
)

func BenchmarkSimpleFetching(b *testing.B) {
	input := []byte(`{"message":"hello world"}`)

	data := map[string]interface{}{}
	json.Unmarshal(input, &data)

	store := NewReferenceStore(len(data))
	store.StoreValues("input", "", data)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Load("input", "message")
	}
}

func BenchmarkSimpleUnmarshal(b *testing.B) {
	input := []byte(`{"message":"hello world"}`)

	data := map[string]interface{}{}
	json.Unmarshal(input, &data)

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewReferenceStore(len(data))
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
		store := NewReferenceStore(len(data))
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
		store := NewReferenceStore(len(data))
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
		store := NewReferenceStore(len(data))
		store.StoreValues("input", "", data)
	}
}

func TestStoreReference(t *testing.T) {
	store := NewReferenceStore(10)

	resource := "input"
	ref := &Reference{
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
	store := NewReferenceStore(1)

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
	store := NewReferenceStore(1)

	nested := "nested"
	key := "message"
	target := template.JoinPath(nested, key)
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

func TestStoreRepeatedMap(t *testing.T) {
	store := NewReferenceStore(1)

	nested := "nested"
	key := "message"
	target := template.JoinPath(nested, key)
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: []map[string]interface{}{
			{
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

func TestStoreRepeatedValues(t *testing.T) {
	store := NewReferenceStore(1)

	nested := "nested"
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: []interface{}{
			value,
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
	result = repeating.Load("", "")
	if result == nil {
		t.Fatal("did not return repeating reference")
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}

func TestStoreRepeatedAppend(t *testing.T) {
	store := NewReferenceStore(1)

	nested := "nested"
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: []interface{}{
			value,
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
	result = repeating.Load("", "")
	if result == nil {
		t.Fatal("did not return repeating reference")
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}

	original := len(result.Repeated)
	result.Append(NewReferenceStore(1))

	if len(result.Repeated) == original {
		t.Fatal("passed store did not get appended")
	}
}

func TestPrefixStoreValue(t *testing.T) {
	store := NewReferenceStore(1)
	resource := "input"
	prefix := "prefix"
	value := "test"
	path := "key"

	pstore := NewPrefixStore(store, resource, prefix)
	pstore.StoreValue("", path, value)

	ref := pstore.Load(resource, template.JoinPath(prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from prefix store")
	}

	if ref.Value != value {
		t.Fatalf("unexpected reference value '%+v', expected '%+v'", ref.Value, value)
	}

	ref = store.Load(resource, template.JoinPath(prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from reference store")
	}

	if ref.Value != value {
		t.Fatalf("unexpected reference value '%+v', expected '%+v'", ref.Value, value)
	}
}

func TestPrefixStoreValues(t *testing.T) {
	store := NewReferenceStore(1)
	resource := "input"
	prefix := "prefix"
	path := "key"

	value := "value"
	values := map[string]interface{}{
		"key": "value",
	}

	pstore := NewPrefixStore(store, resource, prefix)
	pstore.StoreValues("", "", values)

	ref := pstore.Load(resource, template.JoinPath(prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from prefix store")
	}

	if ref.Value != value {
		t.Fatalf("unexpected value '%s', expected '%s'", ref.Value, value)
	}

	ref = store.Load(resource, template.JoinPath(prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from reference store")
	}

	if ref.Value != value {
		t.Fatalf("unexpected reference value '%+v', expected '%+v'", ref.Value, value)
	}
}

func TestPrefixStoreReference(t *testing.T) {
	store := NewReferenceStore(1)
	resource := "input"
	prefix := "prefix"
	path := "key"

	value := NewReference(path)

	pstore := NewPrefixStore(store, resource, prefix)
	pstore.StoreReference("", value)

	ref := pstore.Load(resource, template.JoinPath(prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from prefix store")
	}

	ref = store.Load(resource, template.JoinPath(prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from reference store")
	}
}

func TestMergeReferences(t *testing.T) {
	key := "key"
	resource := "input"
	path := ""

	left := References{}
	right := References{
		key: &specs.PropertyReference{
			Resource: resource,
			Path:     path,
		},
	}

	left.MergeLeft(right)

	if left[key] == nil {
		t.Fatalf("merge failed, expected '%s' to be available", key)
	}

	if left[key].Resource != resource {
		t.Fatalf("unexpected property reference resource '%s', expected '%s'", left[key].Resource, resource)
	}
}

func TestParameterMapReferences(t *testing.T) {
	type expected struct {
		count  int
		params *specs.ParameterMap
	}

	tests := map[string]*expected{
		"header": {
			count: 1,
			params: &specs.ParameterMap{
				Header: specs.Header{
					"key": &specs.Property{
						Reference: &specs.PropertyReference{
							Path: "key",
						},
					},
				},
			},
		},
		"reference": {
			count: 1,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Reference: &specs.PropertyReference{
						Path: "key",
					},
				},
			},
		},
		"nested": {
			count: 2,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Nested: map[string]*specs.Property{
						"first": {
							Reference: &specs.PropertyReference{
								Path: "key",
							},
						},
						"second": {
							Reference: &specs.PropertyReference{
								Path: "else",
							},
						},
					},
				},
			},
		},
		"single": {
			count: 1,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Nested: map[string]*specs.Property{
						"first": {
							Reference: &specs.PropertyReference{
								Path: "key",
							},
						},
						"second": {
							Reference: &specs.PropertyReference{
								Path: "key",
							},
						},
					},
				},
			},
		},
		"empty": {
			count:  0,
			params: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ref := ParameterReferences(test.params)
			t.Log(ref)
			if len(ref) != test.count {
				t.Fatalf("unexpected amount of references %d, exepected %d", len(ref), test.count)
			}
		})
	}
}
