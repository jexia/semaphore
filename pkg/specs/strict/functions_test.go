package strict

import (
	"testing"

	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
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

	custom := functions.Custom{
		"static": func(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
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
			prop := specs.Property{
				Name: "",
				Path: "message",
				Raw:  input,
			}

			err := PrepareFunction(ctx, nil, nil, &prop, make(functions.Stack), custom)
			if err != nil {
				t.Error(err)
			}

			if prop.Reference.Property == nil {
				t.Fatalf("unexpected reference property, reference property not set '%+v'", prop)
			}

			CompareProperties(t, prop, expected)
		})
	}
}

func TestParseUnavailableFunction(t *testing.T) {
	custom := functions.Custom{}
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
		err := PrepareFunction(ctx, nil, nil, prop, make(functions.Stack), custom)
		if err == nil {
			t.Error("unexpected pass")
		}
	}
}
