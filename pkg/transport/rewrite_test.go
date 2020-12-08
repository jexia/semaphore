package transport

import "testing"

func TestNewRewrite(t *testing.T) {
	t.Run("regexp failure", func(t *testing.T) {
		rewrite, err := NewRewrite(`/{]?<(`, `/newPath/$query`)
		if err == nil {
			t.Fatal("error was expected")
		}

		if rewrite != nil {
			t.Fatal("rewrite function was expected to be nil")
		}
	})

	t.Run("malformed template", func(t *testing.T) {
		rewrite, err := NewRewrite(`/`, `/newPath/<query<`)
		if err == nil {
			t.Fatal("error was expected")
		}

		if rewrite != nil {
			t.Fatal("rewrite function was expected to be nil")
		}
	})

	tests := map[string]struct {
		pattern   string
		template  string
		match     bool
		sourceURL string
		expected  string
	}{
		"no match": {
			`/foo/(?P<number>\d+)`,
			`/<number>`,
			false,
			`/foo/bar`,
			`/foo/bar`,
		},
		"rewrite tail": {
			`/oldPath/(?P<tail>.*)`,
			`/newPath/<tail>`,
			true,
			`/oldPath/foo/bar`,
			`/newPath/foo/bar`,
		},
		"swap segments": {
			`/(?P<first>\w+)/(?P<second>\w+)/(?P<tail>.*)`,
			`/prefix/<second>/<first>/<tail>/suffix`,
			true,
			`/foo/bar/baz/42`,
			`/prefix/bar/foo/baz/42/suffix`,
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			rewrite, err := NewRewrite(test.pattern, test.template)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if rewrite == nil {
				t.Fatal("rewrite function was not expected to be nil")
			}

			actual, ok := rewrite(test.sourceURL)
			if !ok && test.match {
				t.Error("no match found")
			} else if ok && !test.match {
				t.Errorf("unexpected match")
			}

			if actual != test.expected {
				t.Errorf("the output URL %q was expected to be %q", actual, test.expected)
			}
		})

	}
}

func TestCompileTemplate(t *testing.T) {
	t.Run("variable is already opened", func(t *testing.T) {
		if _, err := compileTemplate(`<<`); err == nil {
			t.Fatal("error was expected")
		}
	})

	t.Run("variable was not opened", func(t *testing.T) {
		if _, err := compileTemplate(`>`); err == nil {
			t.Fatal("error was expected")
		}
	})

	t.Run("success", func(t *testing.T) {
		var (
			input       = `/<foo>/<bar>`
			expected    = `/$foo/$bar`
			actual, err = compileTemplate(input)
		)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if actual != expected {
			t.Errorf("the compiled template %q was expected to be %q", actual, expected)
		}
	})
}
