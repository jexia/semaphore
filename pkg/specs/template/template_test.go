package template

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func CompareProperties(t *testing.T, left specs.Property, right specs.Property) {
	if left.Scalar == nil || right.Scalar == nil {
		t.Fatalf("scalar not set left (%+v) right (%+v)", left.Scalar, right.Scalar)
	}

	if left.Scalar.Default != right.Scalar.Default {
		t.Errorf("unexpected default '%s', expected '%s'", left.Scalar.Default, right.Scalar.Default)
	}

	if left.Scalar.Type != right.Scalar.Type {
		t.Errorf("unexpected type '%s', expected '%s'", left.Scalar.Type, right.Scalar.Type)
	}

	if left.Scalar.Label != right.Scalar.Label {
		t.Errorf("unexpected label '%s', expected '%s'", left.Scalar.Label, right.Scalar.Label)
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
	name := ""
	path := "message"

	tests := map[string]specs.Property{
		"'prefix'": {
			Name: name,
			Path: path,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "prefix",
				},
			},
		},
		"'edge''": {
			Name: name,
			Path: path,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Label:   labels.Optional,
					Default: "edge'",
				},
			},
		},
		"input:message": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"input:user-id": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "user-id",
			},
		},
		"input.header:Authorization": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "authorization",
			},
		},
		"input.header:User-Id": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "user-id",
			},
		},
		"input.header:": {
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input.header",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			property, err := ParseContent(path, name, input)
			if err != nil {
				t.Fatal(err)
			}

			if property.Path != expected.Path {
				t.Errorf("unexpected path '%s', expected '%s'", property.Path, expected.Path)
			}

			CompareProperties(t, *property, expected)
		})
	}
}

func TestParseReference(t *testing.T) {
	name := ""
	path := "message"

	tests := map[string]specs.Property{
		"input:message": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"input:user-id": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "user-id",
			},
		},
		"input.header:Authorization": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "authorization",
			},
		},
		"input.header:User-Id": {
			Name: name,
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "user-id",
			},
		},
		"input:": {
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
			},
		},
		"input.header:": {
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input.header",
			},
		},
		"input": {
			Path: path,
			Reference: &specs.PropertyReference{
				Resource: "input",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			property, err := ParseReference(path, name, input)
			if err != nil {
				t.Fatal(err)
			}

			if property.Path != expected.Path {
				t.Errorf("unexpected path '%s', expected '%s'", property.Path, expected.Path)
			}

			CompareProperties(t, *property, expected)
		})
	}
}

func TestParseReferenceErr(t *testing.T) {
	name := ""
	path := "message"

	tests := []string{
		"input:..",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := Parse(ctx, path, name, input)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestUnknownReferencePattern(t *testing.T) {
	name := ""
	path := "message"

	tests := []string{
		"input",
		"value",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseContent(path, name, input)
			if err != nil {
				t.Fatalf("unexpected err %s", err)
			}
		})
	}
}

func TestParseReferenceTemplates(t *testing.T) {
	name := ""

	tests := map[string]specs.Property{
		"{{ input:message }}": {
			Path: "message",
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"{{ input.prop:message }}": {
			Path: "message",
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "message",
			},
		},
		"{{ input.prop:user-id }}": {
			Path: "message",
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "user-id",
			},
		},
		"{{ input:user-id }}": {
			Path: "message",
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "user-id",
			},
		},
		"{{ input.prop:message.prop }}": {
			Path: "message.prop",
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "message.prop",
			},
		},
		"{{ input:message.prop }}": {
			Path: "messsage.prop",
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message.prop",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			property, err := Parse(ctx, expected.Path, name, input)
			if err != nil {
				t.Error(err)
			}

			CompareProperties(t, *property, expected)
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
		result := Is(input)
		if result != expected {
			t.Fatalf("unexpected result %+v, expected %+v", result, expected)
		}
	}
}
