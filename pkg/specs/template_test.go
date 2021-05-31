package specs

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

func TestTemplate_Type(t *testing.T) {
	type fields struct {
		Scalar   *Scalar
		Enum     *Enum
		Repeated Repeated
		Message  Message
		OneOf    OneOf
	}
	tests := []struct {
		name   string
		fields fields
		want   types.Type
	}{
		{
			"return scalar type",
			fields{Scalar: &Scalar{Type: types.Int32}},
			types.Int32,
		},
		{
			"return enum",
			fields{Enum: &Enum{}},
			types.Enum,
		},
		{
			"return array",
			fields{Repeated: Repeated{}},
			types.Array,
		},
		{
			"return message",
			fields{Message: Message{}},
			types.Message,
		},
		{
			"return oneOf",
			fields{OneOf: OneOf{}},
			types.OneOf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := Template{
				Scalar:   tt.fields.Scalar,
				Enum:     tt.fields.Enum,
				Repeated: tt.fields.Repeated,
				Message:  tt.fields.Message,
				OneOf:    tt.fields.OneOf,
			}
			if got := template.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func CompareTemplates(t *testing.T, actual, expected Template) {
	if err := actual.Compare(expected); err != nil {
		t.Errorf("unexpected property: %s", err)
	}

	if expected.Scalar != nil {
		if actual.Scalar.Default != expected.Scalar.Default {
			t.Errorf("unexpected default '%s', expected '%s'", actual.Scalar.Default, expected.Scalar.Default)
		}
	}

	if expected.Reference != nil && actual.Reference == nil {
		t.Error("reference not set but expected")
	}

	if expected.Reference != nil {
		if actual.Reference.Resource != expected.Reference.Resource {
			t.Errorf("unexpected reference resource '%s', expected '%s'", actual.Reference.Resource, expected.Reference.Resource)
		}

		if actual.Reference.Path != expected.Reference.Path {
			t.Errorf("unexpected reference path '%s', expected '%s'", actual.Reference.Path, expected.Reference.Path)
		}
	}
}

func TestGetTemplateContent(t *testing.T) {
	tests := map[string]string{
		"{{ input:message }}":            "input:message",
		"{{input:message }}":             "input:message",
		"{{ input:message}}":             "input:message",
		"{{input:message}}":              "input:message",
		"{{input.header:Authorization}}": "input.header:Authorization",
		"{{ add(input:message) }}":       "add(input:message)",
		"{{ add(input:user-id) }}":       "add(input:user-id)",
	}

	for input, expected := range tests {
		result := GetTemplateContent(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestParseTemplateContent(t *testing.T) {
	tests := map[string]Template{
		// Template number (not supported)
		"123": {},
		// Template string
		"'prefix'": {
			Scalar: &Scalar{
				Type:    types.String,
				Default: "prefix",
			},
		},
		// Template string with extra quote
		"'edge''": {
			Scalar: &Scalar{
				Type:    types.String,
				Default: "edge'",
			},
		},
		// Template reference
		"input:message": {
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			property, err := ParseTemplateContent(input)
			if err != nil {
				t.Fatal(err)
			}

			CompareTemplates(t, property, expected)
		})
	}
}

func TestParseReference(t *testing.T) {
	tests := map[string]Template{
		// Reference
		"input:message": {
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		// Error: illegal path
		"input..field": {},
	}

	for input, expected := range tests {
		input := input
		expected := expected

		t.Run(input, func(t *testing.T) {
			property, err := ParseTemplateReference(input)

			if expected.Reference == nil {
				// Could not make the reference due to an error
				if err == nil {
					t.Fatal("Expected an error")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}

				CompareTemplates(t, property, expected)
			}
		})
	}
}

func TestUnknownReferencePattern(t *testing.T) {
	tests := []string{
		"input",
		"value",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseTemplateContent(input)
			if err != nil {
				t.Fatalf("unexpected err %s", err)
			}
		})
	}
}

func TestParseTemplates(t *testing.T) {
	tests := map[string]Template{
		"{{'prefix'}}": {
			Scalar: &Scalar{
				Type:    types.String,
				Default: "prefix",
			},
		},
		"{{ input:message }}": {
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			template, err := ParseTemplate(ctx, "", input)
			if err != nil {
				t.Error(err)
			}

			CompareTemplates(t, template, expected)
		})
	}
}

func TestIsTemplate(t *testing.T) {
	tests := map[string]bool{
		"{{ resource:path }}": true,
		"{{resource:path}}":   true,
		"resource:path":       false,
		"{{ resource:path":    false,
		"{{resource:path":     false,
		"resource:path }}":    false,
		"resource:path}}":     false,
	}

	for input, expected := range tests {
		result := IsTemplate(input)
		if result != expected {
			t.Fatalf("unexpected result %+v, expected %+v", result, expected)
		}
	}
}
