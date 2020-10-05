package prettyerr

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestPrettify(t *testing.T) {
	t.Run("build Errors from several wrapped errors", func(t *testing.T) {
		errOne := errors.New("missing everything")
		errTwo := fmt.Errorf("failed to do One: %w", errOne)

		stack, err := Prettify(errTwo)

		if err != nil {
			t.Errorf("Prettify() is not expected to return an error")
		}

		if len(stack) != 2 {
			t.Errorf("Prettify() is expected to return 3 elements, got: %v", len(stack))
		}

		prettyOne := Error{
			Original: errOne,
			Message:  errOne.Error(),
			Details:  nil,
			Code:     GenericErrorCode,
		}

		prettyTwo := Error{
			Original: errTwo,
			Message:  errTwo.Error(),
			Details:  nil,
			Code:     GenericErrorCode,
		}

		if !reflect.DeepEqual(stack[0], prettyTwo) {
			t.Errorf("Prettify()[0] = %v, want %v", stack[0], prettyTwo)
		}

		if !reflect.DeepEqual(stack[1], prettyOne) {
			t.Errorf("Prettify()[1] = %v, want %v", stack[1], prettyOne)
		}
	})
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Original   error
		Message    string
		Details    map[string]interface{}
		Code       string
		Suggestion string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"returns Message",
			fields{Message: "failed to fail"},
			"failed to fail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Original:   tt.fields.Original,
				Message:    tt.fields.Message,
				Details:    tt.fields.Details,
				Code:       tt.fields.Code,
				Suggestion: tt.fields.Suggestion,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type fields struct {
		Original   error
		Message    string
		Details    map[string]interface{}
		Code       string
		Suggestion string
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			"returns the original error",
			fields{Original: io.EOF},
			io.EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Original:   tt.fields.Original,
				Message:    tt.fields.Message,
				Details:    tt.fields.Details,
				Code:       tt.fields.Code,
				Suggestion: tt.fields.Suggestion,
			}
			if got := e.Unwrap(); got != tt.want {
				t.Errorf("Unwrap() = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func TestStandardErr(t *testing.T) {
	t.Run("build StandardErr from several unwrapped errors", func(t *testing.T) {
		errOne := errors.New("missing everything")
		errTwo := fmt.Errorf("failed to do One: %w", errOne)

		prettyOne := StandardErr(errTwo)

		if reflect.DeepEqual(prettyOne, NoPrettifierErr) {
			t.Errorf("StandardErr() is not expected to return NoPrettifierErr")
		}

	})
}
