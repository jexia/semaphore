package template

import "testing"

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

func TestPathError(t *testing.T) {
	type fields struct {
		Path string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Path: "/semaphore"},
			"unable to resolve path '/semaphore'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrPathNotFound{
				Path: "/semaphore",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
