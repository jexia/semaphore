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

	options := NewOptions(WithPath(path, recursive), WithFunctions(functions), WithSchema(schema))

	if path != options.Path {
		t.Errorf("unexpected path %+v, expected %+v", options.Path, path)
	}

	if recursive != options.Recursive {
		t.Errorf("unexpected recursive definition %+v, expected %+v", options.Recursive, recursive)
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
		{WithPath(".", false), WithSchema(schema)}: true,
		{WithPath(".", false)}:                     true,
		{WithSchema(schema)}:                       false,
		{WithFunctions(functions)}:                 false,
		{WithPath(".", false), WithSchema(schema), WithFunctions(functions)}: true,
	}

	for input, pass := range tests {
		_, err := New(*input...)
		if err == nil && !pass {
			t.Fatalf("unexpected pass %+v, input: %+v", err, input)
		}

		if err != nil && pass {
			t.Fatalf("unexpected fail %+v", err)
		}
	}
}
