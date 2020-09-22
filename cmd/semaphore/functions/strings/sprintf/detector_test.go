package sprintf

import (
	"errors"
	"testing"
)

func TestRadix(t *testing.T) {
	var detector = NewRadix()

	t.Run("should register formatter", func(t *testing.T) {
		if err := detector.Register(String{}); err != nil {
			t.Errorf("unexpected error '%s'", err)
		}
	})

	t.Run("should return an error when conflict is detected", func(t *testing.T) {
		var err = detector.Register(String{})

		if !errors.As(err, &errVerbConflict{}) {
			t.Errorf("unexpected error '%s'", err)
		}
	})

	t.Run("should return 'nil' and 'false' by unknown verb", func(t *testing.T) {
		var constructor, ok = detector.Detect("u")

		if ok {
			t.Error("unexpected 'true'")
		}

		if constructor != nil {
			t.Errorf("unexpected constructor '%T', 'nil' was expected", constructor)
		}
	})

	t.Run("should return available formatter by registered verb", func(t *testing.T) {
		var (
			expected        = String{}
			constructor, ok = detector.Detect("s")
		)

		if !ok {
			t.Error("unexpected 'false'")
		}

		if constructor != expected {
			t.Errorf("unexpected constructor '%T', expected '%T'", constructor, expected)
		}
	})
}
