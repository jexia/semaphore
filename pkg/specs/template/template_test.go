package template

import (
	"testing"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
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

func TestGetTemplateContent(t *testing.T) {
	tests := map[string]string{
		"{{ input:message }}":      "input:message",
		"{{input:message }}":       "input:message",
		"{{ input:message}}":       "input:message",
		"{{input:message}}":        "input:message",
		"{{ add(input:message) }}": "add(input:message)",
	}

	for input, expected := range tests {
		result := GetTemplateContent(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestParseReference(t *testing.T) {
	name := ""
	path := "message"

	tests := map[string]specs.Property{
		"input:message": specs.Property{
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"input:": specs.Property{
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
			},
		},
		"input": specs.Property{
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
			},
		},
	}

	for input, expected := range tests {
		property := ParseReference(path, name, input)

		if property.Path != expected.Path {
			t.Errorf("unexpected path '%s', expected '%s'", property.Path, expected.Path)
		}

		CompareProperties(t, *property, expected)
	}
}

func TestParseTemplate(t *testing.T) {
	name := ""

	tests := map[string]specs.Property{
		"{{ input:message }}": specs.Property{
			Path: "message",
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"{{ input.prop:message }}": specs.Property{
			Path: "message",
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "message",
			},
		},
		"{{ input.prop:message.prop }}": specs.Property{
			Path: "message.prop",
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "message.prop",
			},
		},
		"{{ input:message.prop }}": specs.Property{
			Path: "messsage.prop",
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message.prop",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := instance.NewContext()
			property, err := Parse(ctx, expected.Path, name, input)
			if err != nil {
				t.Error(err)
			}

			CompareProperties(t, *property, expected)
		})
	}
}
