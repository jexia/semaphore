package specs

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
			t.Errorf("unexpected result %v, expected %v", result, expected)
		}

		for index, part := range result {
			if part != expected[index] {
				t.Errorf("unexpected result: %s expected %s", part, expected)
			}
		}
	}
}

func TestPathError(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "return the formatted error",
			path: "/semaphore",
			want: "unable to resolve path '/semaphore'",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			e := ErrPathNotFound{
				Path: tt.path,
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
