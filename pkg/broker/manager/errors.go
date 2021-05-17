package manager

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/v2/pkg/prettyerr"
)

// ErrNilFlowManager is thrown when a nil flow manager has been passed
var ErrNilFlowManager = errors.New("nil flow manager")

// Details contains error details.
type Details struct {
	Service, Flow, Node string
}

// Details returns error details.
func (d Details) Details() map[string]interface{} {
	details := make(map[string]interface{})

	if d.Flow != "" {
		details["Flow"] = d.Flow
	}

	if d.Node != "" {
		details["Node"] = d.Node
	}

	if d.Service != "" {
		details["Service"] = d.Service
	}

	return details
}

// ErrNoServiceName occurs when service name is not provided.
type ErrNoServiceName Details

func (e ErrNoServiceName) Error() string {
	return "no service name configured"
}

// Prettify returns the prettified version of the given error.
func (e ErrNoServiceName) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "EmptyServiceName",
		Message: e.Error(),
		Details: Details(e).Details(),
	}
}

// ErrNoService is invoked when the requested service is unavailable.
type ErrNoService Details

func (e ErrNoService) Error() string {
	return fmt.Sprintf("the service '%s' was not found", e.Service)
}

// Prettify returns the prettified version of the given error.
func (e ErrNoService) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "InvalidService",
		Message: e.Error(),
		Details: Details(e).Details(),
	}
}

// ErrNoTransport indicates that requested transport is unavailable.
type ErrNoTransport struct {
	Details

	Transport string
}

func (e ErrNoTransport) Error() string {
	return fmt.Sprintf("transport '%s' is unavailable", e.Transport)
}

// Prettify returns the prettified version of the given error.
func (e ErrNoTransport) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "InvalidService",
		Message: e.Error(),
		Details: e.Details.Details(),
	}
}

type noCodec struct {
	Details

	Codec string
}

// ErrNoRequestCodec is invoked when request codec is invalid/unknown.
type ErrNoRequestCodec noCodec

func (e ErrNoRequestCodec) Error() string {
	return fmt.Sprintf("request codec '%s' is unavailable", e.Codec)
}

// Prettify returns the prettified version of the given error
func (e ErrNoRequestCodec) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "InvalidCodec",
		Message: e.Error(),
		Details: e.Details.Details(),
	}
}

// ErrNoResponseCodec is invoked when response codec is invalid/unknown.
type ErrNoResponseCodec noCodec

func (e ErrNoResponseCodec) Error() string {
	return fmt.Sprintf("response codec '%s' is unavailable", e.Codec)
}

// Prettify returns the prettified version of the given error
func (e ErrNoResponseCodec) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "InvalidCodec",
		Message: e.Error(),
		Details: e.Details.Details(),
	}
}
