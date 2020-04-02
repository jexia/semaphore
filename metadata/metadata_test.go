package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestMetadataAppend(t *testing.T) {
	md := MD{
		"foo": "bar",
		"bar": "baz",
	}

	append := MD{
		"foo": "bar",
		"bar": "baz",
	}

	md.Append(append)

	for k := range append {
		if _, ok := md[k]; !ok {
			t.Fatalf("Appended metadata not available %s", k)
		}
	}
}

func TestMetadataGet(t *testing.T) {
	expected := "baz"
	md := MD{
		"foo": "bar",
		"bar": expected,
	}

	ctx := NewContext(context.Background(), md)
	result, has := Get(ctx, "bar")
	if !has {
		t.Fatal("key not available inside metadata")
	}

	if result != expected {
		t.Fatalf("unexpected result '%s', expected '%s'", result, expected)
	}
}

func TestMetadataCopy(t *testing.T) {
	md := MD{
		"foo": "bar",
		"bar": "baz",
	}

	cp := Copy(md)

	for k, v := range md {
		if cv := cp[k]; cv != v {
			t.Fatalf("Got %s:%s for %s:%s", k, cv, k, v)
		}
	}
}

func TestMetadataContext(t *testing.T) {
	md := MD{
		"foo": "bar",
	}

	ctx := NewContext(context.TODO(), md)

	emd, ok := FromContext(ctx)
	if !ok {
		t.Errorf("Unexpected error retrieving metadata, got %t", ok)
	}

	if emd["foo"] != md["foo"] {
		t.Errorf("Expected key: %s val: %s, got key: %s val: %s", "foo", md["foo"], "foo", emd["foo"])
	}

	if i := len(emd); i != 1 {
		t.Errorf("Expected metadata length 1 got %d", i)
	}
}

func TestMergeContext(t *testing.T) {
	type args struct {
		existing  MD
		append    MD
		overwrite bool
	}
	tests := []struct {
		name string
		args args
		want MD
	}{
		{
			name: "matching key, overwrite false",
			args: args{
				existing:  MD{"foo": "bar", "sumo": "demo"},
				append:    MD{"sumo": "demo2"},
				overwrite: false,
			},
			want: MD{"foo": "bar", "sumo": "demo"},
		},
		{
			name: "matching key, overwrite true",
			args: args{
				existing:  MD{"foo": "bar", "sumo": "demo"},
				append:    MD{"sumo": "demo2"},
				overwrite: true,
			},
			want: MD{"foo": "bar", "sumo": "demo2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(context.Background(), tt.args.existing)
			merge := MergeContext(ctx, tt.args.append, tt.args.overwrite)

			if got, _ := FromContext(merge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
