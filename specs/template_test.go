package specs

import (
	"testing"

	"github.com/jexia/maestro/internal/instance"
)

func CompareProperties(t *testing.T, left Property, right Property) {
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

func TestJoinPath(t *testing.T) {
	tests := map[string][]string{
		"echo":         {".", "echo"},
		"service.echo": {"service", "echo"},
		"ping.pong":    {"ping.", "pong"},
		"call.me":      {"call.", "me."},
		"":             {"", ""},
		".":            {"", "."},
	}

	for expected, input := range tests {
		result := JoinPath(input...)
		if result != expected {
			t.Errorf("unexpected result: %s expected %s", result, expected)
		}
	}
}

func TestSplitPath(t *testing.T) {
	tests := map[string][]string{
		"service.echo": {"service", "echo"},
		"ping.pong":    {"ping", "pong"},
		"call.me":      {"call", "me"},
		"":             {""},
	}

	for input, expected := range tests {
		result := SplitPath(input)
		if len(result) != len(expected) {
			t.Errorf("unexepcted result %v, expected %v", result, expected)
		}

		for index, part := range result {
			if part != expected[index] {
				t.Errorf("unexpected result: %s expected %s", part, expected)
			}
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

	tests := map[string]Property{
		"input:message": Property{
			Name: name,
			Path: path,
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"input:": Property{
			Path: path,
			Reference: &PropertyReference{
				Resource: "input",
			},
		},
		"input": Property{
			Path: path,
			Reference: &PropertyReference{
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

	tests := map[string]Property{
		"{{ input:message }}": Property{
			Path: "message",
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"{{ input.prop:message }}": Property{
			Path: "message",
			Reference: &PropertyReference{
				Resource: "input.prop",
				Path:     "message",
			},
		},
		"{{ input.prop:message.prop }}": Property{
			Path: "message.prop",
			Reference: &PropertyReference{
				Resource: "input.prop",
				Path:     "message.prop",
			},
		},
		"{{ input:message.prop }}": Property{
			Path: "messsage.prop",
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message.prop",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := instance.NewContext()
			property, err := ParseTemplate(ctx, expected.Path, name, input)
			if err != nil {
				t.Error(err)
			}

			CompareProperties(t, *property, expected)
		})
	}
}
