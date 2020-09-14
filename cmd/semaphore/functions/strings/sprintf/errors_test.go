package sprintf

import (
	"errors"
	"testing"
)

type stringer string

func (s stringer) String() string { return string(s) }

func TestErrVerbConflict(t *testing.T) {
	var err = errVerbConflict{stringer("foo")}

	t.Run("build error message", func(t *testing.T) {
		var (
			expected = `verb "foo" is already in use`
			actual   = err.Error()
		)

		if actual != expected {
			t.Errorf("error '%s' was expected to be '%s'", actual, expected)
		}
	})
}

func TestErrScanFormat(t *testing.T) {
	var (
		inner = errors.New("something went wrong")
		err   = errScanFormat{
			inner:    inner,
			format:   "develop %s not %h",
			position: 15,
		}
	)

	t.Run("build error message", func(t *testing.T) {
		var (
			expected = "something went wrong:\n               ↓\ndevelop %s not %h\n"
			actual   = err.Error()
		)

		if actual != expected {
			t.Errorf("message\n%s\n was expected to be%s", actual, expected)
		}
	})
}
