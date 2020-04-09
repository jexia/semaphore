package strict

import (
	"testing"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
)

func CompareProperties(t *testing.T, left specs.Property, right specs.Property) {
	if left.Default != right.Default {
		t.Errorf("unexpected default '%s', expected '%s'", left.Default, right.Default)
	}

	if left.Type != right.Type {
		t.Errorf("unexpected type '%s', expected '%s'", left.Type, right.Type)
	}

	if left.Label != right.Label {
		t.Errorf("unexpected label '%s', expected '%s'", left.Label, right.Label)
	}

	if right.Reference != nil && left.Reference == nil {
		t.Error("reference not set but expected")
	}

	if right.Reference != nil {
		if left.Reference.Resource != right.Reference.Resource {
			t.Errorf("unexpected reference resource '%s', expected '%s'", left.Reference.Resource, right.Reference.Resource)
		}

		if left.Reference.Path != right.Reference.Path {
			t.Errorf("unexpected reference path '%s', expected '%s'", left.Reference.Path, right.Reference.Path)
		}
	}
}

func TestParseFunction(t *testing.T) {
	static := specs.Property{
		Path:    "message",
		Default: "message",
		Type:    types.String,
		Label:   labels.Optional,
	}

	functions := specs.CustomDefinedFunctions{
		"static": func(args ...*specs.Property) (*specs.Property, specs.FunctionExec, error) {
			return &static, nil, nil
		},
	}

	// NOTE: testing of sub functions is a function specific implementation and is not part of the template library
	tests := map[string]specs.Property{
		"static()": static,
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := instance.NewContext()
			prop := &specs.Property{
				Name: "",
				Path: "message",
				Raw:  input,
			}

			result, err := PrepareFunction(ctx, nil, nil, prop, make(specs.Functions), functions)
			if err != nil {
				t.Error(err)
			}

			if result.Reference == nil {
				t.Fatalf("unexpected property reference, reference not set '%+v'", result)
			}

			if result.Reference.Property == nil {
				t.Fatalf("unexpected reference property, reference property not set '%+v'", result)
			}

			CompareProperties(t, *result.Reference.Property, expected)
		})
	}
}

func TestParseUnavailableFunction(t *testing.T) {

	functions := specs.CustomDefinedFunctions{}

	tests := []string{
		"add()",
	}

	for _, input := range tests {
		prop := &specs.Property{
			Name: "",
			Path: "message",
			Raw:  input,
		}

		ctx := instance.NewContext()
		_, err := PrepareFunction(ctx, nil, nil, prop, make(specs.Functions), functions)
		if err == nil {
			t.Error("unexpected pass")
		}
	}
}
