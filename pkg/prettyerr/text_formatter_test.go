package prettyerr

import "testing"

func TestTextFormatter(t *testing.T) {
	missingConfig := Error{
		Message: "File config is missing.",
		Details: map[string]interface{}{
			"filename": "/etc/foo/config.yml",
		},
		Code: "MissingConfig",
	}

	genericError := Error{
		Message: "Something happened, I'm sorry.",
		Details: nil,
		Code:    "GenericError",
	}

	want :=
		"File config is missing.\n" +
			"\tfilename: /etc/foo/config.yml\n" +
			"Something happened, I'm sorry.\n"

	stack := Errors{missingConfig, genericError}
	format := "{{ .Message }}\n{{ range $key, $value := .Details }}\t{{ $key }}: {{ $value }}\n{{ end }}"

	got, err := TextFormatter(stack, format)
	if err != nil {
		t.Errorf("TextFormatter() returned error: %v", err)
	}

	if got != want {
		t.Errorf("TextFormatter(). Got:\n%v\nWant:\n%v", got, want)
	}
}
