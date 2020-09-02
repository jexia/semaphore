package pretty_errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNewStack(t *testing.T) {
	t.Run("build Stack from several wrapped errors", func(t *testing.T) {
		errOne := errors.New("missing everything")
		errTwo := fmt.Errorf("failed to do One: %w", errOne)

		stack, err := NewStack(errTwo)

		if err != nil {
			t.Errorf("NewStack() is not expected to return an error")
		}

		if len(stack) != 2 {
			t.Errorf("NewStack() is expected to return 3 elements, got: %v", len(stack))
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
			t.Errorf("NewStack()[0] = %v, want %v", stack[0], prettyTwo)
		}

		if !reflect.DeepEqual(stack[1], prettyOne) {
			t.Errorf("NewStack()[1] = %v, want %v", stack[1], prettyOne)
		}
	})
}
