package manager

import (
	"reflect"
	"testing"

	"github.com/jexia/semaphore/v2/pkg/prettyerr"
)

func checkMessage(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("the error message %q was expected to be %q", actual, expected)
	}
}

func checkPretty(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("prettyfied error:\n%#v\n does not match the expected one:\n%#v", actual, expected)
	}
}

func TestDetails(t *testing.T) {
	var (
		details = Details{
			Flow:    "flow",
			Node:    "node",
			Service: "service",
		}

		expected = map[string]interface{}{
			"Flow":    "flow",
			"Node":    "node",
			"Service": "service",
		}
	)

	checkPretty(t, details.Details(), expected)
}

func TestErrNoServiceName(t *testing.T) {
	var (
		err = ErrNoServiceName{
			Flow: "flow",
			Node: "node",
		}

		message = "no service name configured"

		expected = prettyerr.Error{
			Code:    "EmptyServiceName",
			Message: message,
			Details: map[string]interface{}{
				"Flow": "flow",
				"Node": "node",
			},
		}
	)

	checkMessage(t, err.Error(), message)
	checkPretty(t, err.Prettify(), expected)
}

func TestErrNoService(t *testing.T) {
	var (
		err = ErrNoService{
			Flow:    "flow",
			Node:    "node",
			Service: "service",
		}

		message = "the service 'service' was not found"

		expected = prettyerr.Error{
			Code:    "InvalidService",
			Message: message,
			Details: map[string]interface{}{
				"Flow":    "flow",
				"Node":    "node",
				"Service": "service",
			},
		}
	)

	checkMessage(t, err.Error(), message)
	checkPretty(t, err.Prettify(), expected)
}

func TestErrNoTransport(t *testing.T) {
	var (
		err = ErrNoTransport{
			Transport: "grpc",
			Details: Details{
				Flow:    "flow",
				Node:    "node",
				Service: "service",
			},
		}

		message = "transport 'grpc' is unavailable"

		expected = prettyerr.Error{
			Code:    "InvalidService",
			Message: message,
			Details: map[string]interface{}{
				"Flow":    "flow",
				"Node":    "node",
				"Service": "service",
			},
		}
	)

	checkMessage(t, err.Error(), message)
	checkPretty(t, err.Prettify(), expected)
}

func TestErrNoRequestCodec(t *testing.T) {
	var (
		err = ErrNoRequestCodec{
			Codec: "json",
			Details: Details{
				Flow:    "flow",
				Node:    "node",
				Service: "service",
			},
		}

		message = "request codec 'json' is unavailable"

		expected = prettyerr.Error{
			Code:    "InvalidCodec",
			Message: message,
			Details: map[string]interface{}{
				"Flow":    "flow",
				"Node":    "node",
				"Service": "service",
			},
		}
	)

	checkMessage(t, err.Error(), message)
	checkPretty(t, err.Prettify(), expected)
}

func TestErrNoResponseCodec(t *testing.T) {
	var (
		err = ErrNoResponseCodec{
			Codec: "json",
			Details: Details{
				Flow:    "flow",
				Node:    "node",
				Service: "service",
			},
		}

		message = "response codec 'json' is unavailable"

		expected = prettyerr.Error{
			Code:    "InvalidCodec",
			Message: message,
			Details: map[string]interface{}{
				"Flow":    "flow",
				"Node":    "node",
				"Service": "service",
			},
		}
	)

	checkMessage(t, err.Error(), message)
	checkPretty(t, err.Prettify(), expected)
}
