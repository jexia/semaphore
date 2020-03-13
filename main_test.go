package maestro

import (
	"testing"

	"github.com/jexia/maestro/specs"
)

func TestOptions(t *testing.T) {
	functions := map[string]specs.PrepareCustomFunction{
		"cdf": nil,
	}

	options := NewOptions(WithFunctions(functions))

	if len(options.Functions) != len(functions) {
		t.Errorf("unexpected functions %+v, expected %+v", options.Functions, functions)
	}
}

func TestNew(t *testing.T) {
	functions := map[string]specs.PrepareCustomFunction{
		"cdf": nil,
	}

	tests := [][]Option{
		{WithDefinitions(nil), WithSchema(nil)},
		{WithDefinitions(nil)},
		{WithSchema(nil)},
		{WithFunctions(functions)},
		{WithDefinitions(nil), WithSchema(nil), WithFunctions(functions)},
	}

	for _, input := range tests {
		_, err := New(input...)
		if err != nil {
			t.Fatalf("unexpected fail %+v", err)
		}
	}
}
