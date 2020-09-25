package references

import (
	"testing"

	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/specs"
)

type SomeExpression struct {
}

func (SomeExpression) Position() string {
	return "7:40"
}

func TestErrUndefinedReference_Prettify(t *testing.T) {
	t.Run("includes Expression in details", func(t *testing.T) {
		expr := SomeExpression{}
		err := ErrUndefinedReference{
			Property: &specs.Property{
				Expr: expr,
			},
		}

		pretty := err.Prettify().Details
		got := pretty["Expression"]
		want := expr.Position()

		if got != want {
			t.Errorf("Details[Expression] = %v, expected %v", got, want)
		}
	})

	t.Run("does not include nil Expression in details", func(t *testing.T) {
		err := ErrUndefinedReference{
			Property: &specs.Property{},
		}

		pretty := err.Prettify().Details
		got, ok := pretty["Expression"]

		if ok {
			t.Errorf("Details[Expression] = %v, expected not being included", got)
		}
	})

	t.Run("returns pretty error", func(t *testing.T) {
		err := ErrUndefinedReference{
			Property: &specs.Property{
				Path: "there",
				Template: specs.Template{
					Reference: &specs.PropertyReference{Resource: "user", Path: "name"},
				},
			},
			Breakpoint: "here",
		}

		got := err.Prettify()

		want := prettyerr.Error{
			Message: err.Error(),
			Details: map[string]interface{}{
				"Reference":  err.Property.Reference,
				"Breakpoint": err.Breakpoint,
				"Path":       err.Property.Path,
			},
			Code:       "UndefinedReference",
			Suggestion: "",
		}

		if got.Error() != want.Error() {
			t.Fatalf("Pretty() = %#v, want %#v", got, want)
		}
	})
}
