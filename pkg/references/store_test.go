package references

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/jexia/semaphore/v2/pkg/specs"
)

func BenchmarkSimpleFetching(b *testing.B) {
	store := NewStore(1)
	store.Store("input:message", &Reference{Value: "hello world"})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Load("input:message")
	}
}

func BenchmarkStoreSingleValue(b *testing.B) {
	input := "hello world"
	store := NewStore(b.N)

	keys := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Store("input:"+keys[i], &Reference{Value: input})
	}
}

func BenchmarkStoreSingleEnum(b *testing.B) {
	enum := int32(1)
	store := NewStore(b.N)

	keys := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Store("input:"+keys[i], &Reference{Enum: &enum})
	}
}

func BenchmarkSimpleUnmarshal(b *testing.B) {
	input := []byte(`{"message":"hello world"}`)

	data := map[string]interface{}{}
	err := json.Unmarshal(input, &data)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewStore(len(data))
		store.Store("input:", &Reference{Value: data})
	}
}

func BenchmarkRepeatedUnmarshal(b *testing.B) {
	input := []byte(`{"details":[{"message":"hello world"}]}`)

	data := map[string]interface{}{}
	err := json.Unmarshal(input, &data)
	if err != nil {
		b.Fatal(err)
	}

	store := NewStore(len(data) * b.N)

	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Store("input:", &Reference{Value: data})
	}
}

func TestEnumJSON(t *testing.T) {
	expected := "PENDING"
	val := Enum(expected, 1)
	bb, err := val.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	want := strconv.Quote(expected)
	if string(bb) != want {
		t.Fatalf("unexpected result '%s', expected '%s'", string(bb), want)
	}

	err = val.UnmarshalJSON([]byte(want))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStoreReference(t *testing.T) {
	store := NewStore(10)

	path := specs.ResourcePath(specs.InputResource, "key")
	ref := &Reference{
		Value: "hello world",
	}

	store.Store(path, ref)
	result := store.Load(path)
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Value != ref.Value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, ref.Value)
	}
}

func TestStoreLoadUnkownReference(t *testing.T) {
	store := NewStore(0)
	result := store.Load("")
	if result != nil {
		t.Fatal("unexpected result")
	}
}

func TestStoreValues(t *testing.T) {
	store := NewStore(1)
	tracker := NewTracker()

	target := "message"
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		target: value,
	}

	StoreValues(store, tracker, specs.ResourcePath(resource), values)
	result := store.Load(specs.ResourcePath(resource, target))
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}

func TestStoreNestedValues(t *testing.T) {
	store := NewStore(1)
	tracker := NewTracker()

	nested := "nested"
	key := "message"
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: map[string]interface{}{
			key: value,
		},
	}

	StoreValues(store, tracker, specs.ResourcePath(resource), values)

	result := store.Load(specs.ResourcePath(resource, nested, key))
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}

func TestStoreRepeatedMap(t *testing.T) {
	store := NewStore(1)
	tracker := NewTracker()

	nested := "nested"
	key := "message"
	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		nested: []map[string]interface{}{
			{
				key: value,
			},
		},
	}

	StoreValues(store, tracker, specs.ResourcePath(resource), values)

	length := store.Length(specs.ResourcePath(resource, nested))
	if length != 1 {
		t.Fatalf("unexpected repeated length %d, expected 1", length)
	}

	arrayp := specs.ResourcePath(resource, nested)
	path := specs.ResourcePath(resource, nested, key)
	tracker.Track(arrayp, 0)

	for index := 0; index < length; index++ {
		result := store.Load(tracker.Resolve(path))
		if result == nil {
			t.Fatal("did not return repeating reference")
		}

		if result.Value != value {
			t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
		}

		tracker.Next(arrayp)
	}
}

func TestStoreEnum(t *testing.T) {
	store := NewStore(1)
	tracker := NewTracker()

	key := "enum"
	expected := int32(1)

	resource := "input"
	values := map[string]interface{}{
		key: Enum("PENDING", expected),
	}

	StoreValues(store, tracker, specs.ResourcePath(resource), values)

	path := specs.ResourcePath(resource, key)
	result := store.Load(path)
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Enum == nil {
		t.Fatal("unexpected result expected enum to be set")
	}

	if *result.Enum != 1 {
		t.Fatalf("unexpected enum result '%d', expected '%d'", result.Enum, expected)
	}
}

func TestStoreRepeatingEnum(t *testing.T) {
	store := NewStore(1)
	tracker := NewTracker()

	key := "enum"
	expected := []interface{}{
		Enum("PENDING", 1),
		Enum("UNKNOWN", 0),
	}

	resource := "input"
	values := map[string]interface{}{
		key: expected,
	}

	StoreValues(store, tracker, specs.ResourcePath(resource), values)

	path := specs.ResourcePath(resource, key)
	length := store.Length(path)
	if length != len(expected) {
		t.Fatalf("unexpected repeated store length %d, expected %d", length, len(expected))
	}

	tracker.Track(path, 0)
	for index := 0; index < length; index++ {
		ref := store.Load(tracker.Resolve(path))
		if ref == nil {
			t.Fatal("unexpected empty reference expected reference to be returned")
		}

		if ref.Enum == nil {
			t.Fatal("unexpected enum expected enum to be defined")
		}

		want := expected[index].(*EnumVal).pos
		if *ref.Enum != want {
			t.Fatalf("unexpected enum %d, expected %d", ref.Enum, want)
		}

		tracker.Next(path)
	}
}

func TestStoreRepeatedValues(t *testing.T) {
	store := NewStore(1)
	tracker := NewTracker()

	value := "hello world"

	resource := "input"
	values := map[string]interface{}{
		"nested": []interface{}{
			value,
		},
	}

	StoreValues(store, tracker, specs.ResourcePath(resource), values)

	length := store.Length(specs.ResourcePath(resource, "nested"))
	if length != 1 {
		t.Errorf("unexpected length %d, expected 1", length)
	}

	result := store.Load(specs.ResourcePath(resource, "nested[0]"))
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Value != value {
		t.Fatalf("unexpected value %+v, expected %+v", result.Value, value)
	}
}

func TestPrefixStoreEnum(t *testing.T) {
	store := NewStore(1)

	expected := int32(1)
	prefix := "prefix"
	path := "key"

	resource := "input"

	pstore := NewPrefixStore(store, specs.ResourcePath(resource, prefix))
	pstore.Store(path, &Reference{Enum: &expected})

	result := store.Load(specs.ResourcePath(resource, prefix, path))
	if result == nil {
		t.Fatal("did not return reference")
	}

	if result.Enum == nil {
		t.Fatal("unexpected result expected enum to be set")
	}

	if *result.Enum != expected {
		t.Fatalf("unexpected enum result '%d', expected '%d'", result.Enum, expected)
	}
}

func TestPrefixStoreValues(t *testing.T) {
	store := NewStore(1)
	resource := "input"
	prefix := "prefix"
	path := "key"

	value := "value"
	values := map[string]interface{}{
		"key": "value",
	}

	pstore := NewPrefixStore(store, specs.ResourcePath(resource, prefix))
	tracker := NewTracker()

	StoreValues(pstore, tracker, "", values)

	ref := store.Load(specs.ResourcePath(resource, prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from prefix store")
	}

	if ref.Value != value {
		t.Fatalf("unexpected value '%s', expected '%s'", ref.Value, value)
	}

	ref = store.Load(specs.ResourcePath(resource, prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from reference store")
	}

	if ref.Value != value {
		t.Fatalf("unexpected reference value '%+v', expected '%+v'", ref.Value, value)
	}
}

func TestPrefixStoreValue(t *testing.T) {
	store := NewStore(1)
	resource := "input"
	prefix := "prefix"
	path := "key"
	value := "message"

	pstore := NewPrefixStore(store, specs.ResourcePath(resource, prefix))
	pstore.Store(path, &Reference{Value: value})

	ref := store.Load(specs.ResourcePath(resource, prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from prefix store")
	}

	ref = store.Load(specs.ResourcePath(resource, prefix, path))
	if ref == nil {
		t.Fatal("unable to load reference from reference store")
	}
}

func TestMergeReferences(t *testing.T) {
	key := "key"
	resource := "input"
	path := ""

	left := Collection{}
	right := Collection{
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
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Path: "key",
							},
						},
					},
				},
			},
		},
		"reference": {
			count: 1,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Path: "key",
						},
					},
				},
			},
		},
		"nested": {
			count: 2,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Template: specs.Template{
						Message: specs.Message{
							"first": {
								Name: "first",
								Path: "first",
								Template: specs.Template{
									Reference: &specs.PropertyReference{
										Path: "key",
									},
								},
							},
							"second": {
								Name: "second",
								Path: "second",
								Template: specs.Template{
									Reference: &specs.PropertyReference{
										Path: "else",
									},
								},
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
					Template: specs.Template{
						Message: specs.Message{
							"first": {
								Name: "first",
								Path: "first",
								Template: specs.Template{
									Reference: &specs.PropertyReference{
										Path: "key",
									},
								},
							},
							"second": {
								Name: "second",
								Path: "second",
								Template: specs.Template{
									Reference: &specs.PropertyReference{
										Path: "key",
									},
								},
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
			if len(ref) != test.count {
				t.Fatalf("unexpected amount of references %d, expected %d", len(ref), test.count)
			}
		})
	}
}
