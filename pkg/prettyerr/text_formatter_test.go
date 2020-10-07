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

	stack := Errors{missingConfig, genericError}

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{
			name:   "No TextFormatter error",
			format: "{{ .Message }}\n{{ range $key, $value := .Details }}\t{{ $key }}: {{ $value }}\n{{ end }}",
			want: "File config is missing.\n" +
				"\tfilename: /etc/foo/config.yml\n" +
				"Something happened, I'm sorry.\n",
		},
		{
			name:   "Fields error",
			format: "{{ .Messag }}\n{{ range $key, $value := .Details }}\t{{ $key }}: {{ $value }}\n{{ end }}",
			want: "failed to execute template for 0: template: node:1:3: executing \"node\" at <.Messag>: can't evaluate field Messag in type prettyerr.Error\n" +
				"File config is missing.\n" +
				"\tfilename: \"/etc/foo/config.yml\"\n" +
				"failed to execute template for 1: template: node:1:3: executing \"node\" at <.Messag>: can't evaluate field Messag in type prettyerr.Error\n" +
				"Something happened, I'm sorry.\n",
		},
		{
			name:   "Template Parsing error",
			format: "{{ .Message }}\n{{ range $key, $value := .Details }}\t{{ $ke }}: {{ $value }}\n{{ end }}",
			want: "failed to parse template: template: node:2: undefined variable \"$ke\"\n" +
				"File config is missing.\n" +
				"\tfilename: \"/etc/foo/config.yml\"\n" +
				"Something happened, I'm sorry.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TextFormatter(stack, tt.format)

			if err != nil {
				t.Errorf("TextFormatter() returned error: %v", err)
			}

			if got != tt.want {
				t.Errorf("TextFormatter(). Got:\n%v\nWant:\n%v", got, tt.want)
			}
		})
	}

}
