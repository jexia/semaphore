package specs

import "testing"

func CompareProperties(t *testing.T, left Property, right Property) {
	if left.Path != right.Path {
		t.Errorf("unexpected path %s, expected %s", left.Path, right.Path)
	}

	if left.Default != right.Default {
		t.Errorf("unexpected default %s, expected %s", left.Default, right.Default)
	}

	if left.Type != right.Type {
		t.Errorf("unexpected type %s, expected %s", left.Type, right.Type)
	}

	if right.Reference != nil && left.Reference == nil {
		t.Error("reference not set but expected")
	}

	if right.Reference != nil {
		if left.Reference.Resource != right.Reference.Resource {
			t.Errorf("unexpected reference resource %s, expected %s", left.Reference.Resource, right.Reference.Resource)
		}

		if left.Reference.Path != right.Reference.Path {
			t.Errorf("unexpected reference path %s, expected %s", left.Reference.Path, right.Reference.Path)
		}
	}

	if right.Function != nil && left.Function == nil {
		t.Error("function not set but expected")
	}
}

func TestJoinPath(t *testing.T) {
	tests := map[string][]string{
		"service.echo": {"service", "echo"},
		"ping.pong":    {"ping.", "pong"},
		"call.me":      {"call.", "me."},
		"":             {"", ""},
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
	path := "message"
	tests := map[string]Property{
		"input:message": Property{
			Path: path,
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
				Label:    LabelOptional,
			},
		},
		"input:": Property{
			Path: path,
			Reference: &PropertyReference{
				Resource: "input",
				Label:    LabelOptional,
			},
		},
		"input": Property{
			Path: path,
			Reference: &PropertyReference{
				Resource: "input",
				Label:    LabelOptional,
			},
		},
	}

	for input, expected := range tests {
		property := ParseReference(path, input)
		CompareProperties(t, *property, expected)
	}
}

func TestParseFunction(t *testing.T) {
	path := "message"
	static := Property{
		Path:    path,
		Default: "message",
		Type:    TypeString,
	}

	functions := CustomDefinedFunctions{
		"add": func(path string, args ...*Property) (*Property, error) {
			return &static, nil
		},
	}

	// NOTE: testing of sub functions is a function specific implementation and is not part of the template libary
	tests := map[string]Property{
		"add(input:message)": static,
	}

	for input, expected := range tests {
		property, err := ParseFunction(path, functions, input)
		if err != nil {
			t.Error(err)
		}

		CompareProperties(t, *property, expected)
	}
}

func TestParseTemplate(t *testing.T) {
	path := "message"
	static := Property{
		Path:    path,
		Default: "message",
		Type:    TypeString,
	}

	functions := CustomDefinedFunctions{
		"add": func(path string, args ...*Property) (*Property, error) {
			return &static, nil
		},
	}

	tests := map[string]Property{
		"{{ input:message }}": Property{
			Path: path,
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		"{{ add(input:message) }}": static,
	}

	for input, expected := range tests {
		property, err := ParseTemplate(path, functions, input)
		if err != nil {
			t.Error(err)
		}

		CompareProperties(t, *property, expected)
	}
}
