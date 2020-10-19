package sprintf

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type stringer string

func (s stringer) String() string { return string(s) }

func TestErrInvalidArguments(t *testing.T) {
	var (
		expected = `invalid number of input arguments 1, expected 2`
		err      = errInvalidArguments{actual: 1, expected: 2}
	)

	t.Run("build error message", func(t *testing.T) {
		var actual = err.Error()
		if actual != expected {
			t.Errorf("error '%s' was expected to be '%s'", actual, expected)
		}
	})
}

func TestErrCannotFormat(t *testing.T) {
	var (
		expected = `cannot use '%foo' formatter for argument 'bar' of type 'bool'`
		err      = errCannotFormat{
			formatter: stringer("foo"),
			argument: &specs.Property{
				Name: "bar",
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.Bool,
					},
				},
			},
		}
	)

	t.Run("build error message", func(t *testing.T) {
		var actual = err.Error()
		if actual != expected {
			t.Errorf("error '%s' was expected to be '%s'", actual, expected)
		}
	})
}

func TestErrVerbConflict(t *testing.T) {
	var (
		expected = `verb "foo" is already in use`
		err      = errVerbConflict{stringer("foo")}
	)

	t.Run("build error message", func(t *testing.T) {
		var actual = err.Error()
		if actual != expected {
			t.Errorf("error '%s' was expected to be '%s'", actual, expected)
		}
	})
}

func TestErrScanFormat(t *testing.T) {
	var (
		expected = "something went wrong:\n               ↓\ndevelop %s not %h\n"
		inner    = errors.New("something went wrong")
		err      = errScanFormat{
			inner:    inner,
			format:   "develop %s not %h",
			position: 15,
		}
	)

	t.Run("build error message", func(t *testing.T) {
		var actual = err.Error()
		if actual != expected {
			t.Errorf("message\n%s\n was expected to be%s", actual, expected)
		}
	})
}
