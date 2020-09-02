package pretty_errors

import (
	"errors"
	"reflect"
	"testing"
)

type errTooPretty struct {
}

func (e errTooPretty) Error() string { return "too pretty" }
func (e errTooPretty) Prettify() Error {
	return Error{
		Original: e,
		Message:  "too pretty. Expected: less pretty",
		Details: nil,
		Code:     "TooPretty",
	}
}

func TestPrettifierStrategy_Match(t *testing.T) {
	strategy := PrettifierStrategy{}

	t.Run("returns the defined prettifier", func(t *testing.T) {
		prettifier := strategy.Match(errTooPretty{})

		got := prettifier.Prettify()
		want := Error{
			Original: errTooPretty{},
			Message:  "too pretty. Expected: less pretty",
			Details: nil,
			Code:     "TooPretty",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Match() is expected to the error's prettifier be returned")
		}
	})

	t.Run("returns a generic prettifier", func(t *testing.T) {
		err := errors.New("fail")
		prettifier := strategy.Match(err)

		got := prettifier.Prettify()
		want := Error{
			Original: err,
			Message:  err.Error(),
			Details: nil,
			Code:     GenericErrorCode,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Match() is expected to the generic prettifier be returned")
		}
	})
}