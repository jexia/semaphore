package maestro

import (
	"testing"

	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs"
)

func TestOptions(t *testing.T) {
	path := "/root"
	recursive := true
	schema := &mock.Collection{}
	functions := map[string]specs.PrepareCustomFunction{
		"cdf": nil,
	}

	options := NewOptions(WithPath(path, recursive), WithFunctions(functions), WithSchemaCollection(schema))

	if path != options.Path {
		t.Errorf("unexpected path %+v, expected %+v", options.Path, path)
	}

	if recursive != options.Recursive {
		t.Errorf("unexpected recursive definition %+v, expected %+v", options.Recursive, recursive)
	}

	if schema != options.Schema {
		t.Errorf("unexpected schema %+v, expected %+v", options.Schema, schema)
	}

	if len(options.Functions) != len(functions) {
		t.Errorf("unexpected functions %+v, expected %+v", options.Functions, functions)
	}
}

func TestNew(t *testing.T) {
	schema := &mock.Collection{}
	functions := map[string]specs.PrepareCustomFunction{
		"cdf": nil,
	}

	tests := map[*[]Option]bool{
		{WithPath(".", false), WithSchemaCollection(schema)}: true,
		{WithSchemaCollection(schema)}:                       false,
		{WithSchemaCollection(schema)}:                       false,
		{WithPath(".", false)}:                               false,
		{WithFunctions(functions)}:                           false,
		{WithPath(".", false), WithSchemaCollection(schema), WithFunctions(functions)}: true,
	}

	for input, pass := range tests {
		_, err := New(*input...)
		if err == nil && !pass {
			t.Fatalf("unexpected pass %s", err)
		}

		if err != nil && pass {
			t.Fatalf("unexpected fail %s", err)
		}
	}
}
