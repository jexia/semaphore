package refs

import (
	"encoding/json"
	"testing"
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
