package template

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

func CompareTemplates(t *testing.T, actual, expected specs.Template) {
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
	tests := map[string]specs.Template{
		"'prefix'": {
			Scalar: &specs.Scalar{
				Type:    types.String,
				Default: "prefix",
			},
		},
		"'edge''": {
			Scalar: &specs.Scalar{
				Type:    types.String,
				Default: "edge'",
			},
		},
		"input:message": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"input:user-id": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "user-id",
			},
		},
		"input.header:Authorization": {
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "authorization",
			},
		},
		"input.header:User-Id": {
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "user-id",
			},
		},
		"input.header:": {
			Reference: &specs.PropertyReference{
				Resource: "input.header",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			property, err := ParseContent(input)
			if err != nil {
				t.Fatal(err)
			}

			CompareTemplates(t, property, expected)
		})
	}
}

func TestParseReference(t *testing.T) {
	tests := map[string]specs.Template{
		"input:message": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"input:user-id": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "user-id",
			},
		},
		"input.header:Authorization": {
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "authorization",
			},
		},
		"input.header:User-Id": {
			Reference: &specs.PropertyReference{
				Resource: "input.header",
				Path:     "user-id",
			},
		},
		"input:": {
			Reference: &specs.PropertyReference{
				Resource: "input",
			},
		},
		"input.header:": {
			Reference: &specs.PropertyReference{
				Resource: "input.header",
			},
		},
		"input": {
			Reference: &specs.PropertyReference{
				Resource: "input",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			property, err := ParseReference(input)
			if err != nil {
				t.Fatal(err)
			}

			CompareTemplates(t, property, expected)
		})
	}
}

func TestParseReferenceErr(t *testing.T) {
	var (
		path  = "message"
		tests = []string{
			"input:..",
		}
	)

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := Parse(ctx, path, input)
			if err == nil {
				t.Fatal("unexpected pass")
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
			_, err := ParseContent(input)
			if err != nil {
				t.Fatalf("unexpected err %s", err)
			}
		})
	}
}

func TestParseReferenceTemplates(t *testing.T) {
	tests := map[string]specs.Template{
		"{{ input:message }}": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"{{ input.prop:message }}": {
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "message",
			},
		},
		"{{ input.prop:user-id }}": {
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "user-id",
			},
		},
		"{{ input:user-id }}": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "user-id",
			},
		},
		"{{ input.prop:message.prop }}": {
			Reference: &specs.PropertyReference{
				Resource: "input.prop",
				Path:     "message.prop",
			},
		},
		"{{ input:message.prop }}": {
			Reference: &specs.PropertyReference{
				Resource: "input",
				Path:     "message.prop",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			template, err := Parse(ctx, "", input)
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
		result := Is(input)
		if result != expected {
			t.Fatalf("unexpected result %+v, expected %+v", result, expected)
		}
	}
}
