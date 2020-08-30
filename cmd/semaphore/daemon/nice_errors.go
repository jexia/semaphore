package daemon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
)

const (
	UnresolvedCallCode         = "UnresolvedCall"
	unresolvedParamsCode       = "UnresolvedParams"
	unresolvedPropertyCode     = "UnresolvedProperty"
	unresolvedParameterMapCode = "UnresolvedParameterMap"
	unableToResolveFlowCode    = "UnableToResolveFlow"
	undefinedReferenceCode     = "UndefinedReference"
	undefinedResourceCode      = "UndefinedResource"
	unresolvedOnErrorCode      = "UnresolvedOnError"
	genericErrorCode           = "GenericError"
	unresolvedNodeCode         = "UnresolvedNode"
)

type NiceError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// NewNiceError takes an internal error and returns a new one that is much nicer to users.
func NewNiceError(err error) error {
	stack := humanizeErrors(err)

	return errors.New(formatNiceErrors(stack))
}

// format stack of Nice Errors with a nice format. For example, it might be a simple json.Marshal call,
// but here we do some friendly magic.
func formatNiceErrors(stack []NiceError) string {
	o := bytes.NewBufferString("\n")

	for _, nice := range stack {
		fmt.Fprintf(o, "(%s) %s\n", nice.Code, nice.Message)
		if nice.Details != nil {
			for k, v := range nice.Details {
				value, _ := json.Marshal(v)
				fmt.Fprintf(o, "\t%s: %s\n", k, value)
			}
		}
	}

	return o.String()
}

func humanizeErrors(next error) []NiceError {
	var (
		stack []NiceError
	)

	for next != nil {
		stack = append(stack, humanizeError(next))
		next = errors.Unwrap(next)
	}

	return stack
}

// check the error and build a friendly error message
func humanizeError(inner error) (nice NiceError) {
	var (
		errUnresolvedFlow         references.ErrUnresolvedFlow
		errUndefinedReference     references.ErrUndefinedReference
		errUndefinedResource      references.ErrUndefinedResource
		errUnresolvedOnError      references.ErrUnresolvedOnError
		errUnresolvedNode         references.ErrUnresolvedNode
		errUnresolvedParameterMap references.ErrUnresolvedParameterMap
		errUnresolvedProperty     references.ErrUnresolvedProperty
		errUnresolvedParams       references.ErrUnresolvedParams
		errUnresolvedCall         references.ErrUnresolvedCall
	)

	switch {
	case errors.As(inner, &errUnresolvedFlow):
		return humanizeErrUnresolvedFlow(errUnresolvedFlow)

	case errors.As(inner, &errUndefinedReference):
		return humanizeErrUndefinedReference(errUndefinedReference)

	case errors.As(inner, &errUndefinedResource):
		return humanizeErrUndefinedResource(errUndefinedResource)

	case errors.As(inner, &errUnresolvedOnError):
		return humanizeErrUnresolvedOnError(errUnresolvedOnError)


	case errors.As(inner, &errUnresolvedNode):
		return humanizeErrUnresolvedNode(errUnresolvedNode)

	case errors.As(inner, &errUnresolvedParameterMap):
		return humanizeErrUnresolvedParameterMap(errUnresolvedParameterMap)


	case errors.As(inner, &errUnresolvedProperty):
		return humanizeErrUnresolvedProperty(errUnresolvedProperty)

	case errors.As(inner, &errUnresolvedParams):
		return humanizeErrUnresolvedParams(errUnresolvedParams)

	case errors.As(inner, &errUnresolvedCall):
		return humanizeErrUnresolvedCall(errUnresolvedCall)

	default:
		return NiceError{
			Code:    genericErrorCode,
			Message: inner.Error(),
		}
	}
}

func humanizeErrUnresolvedFlow(err references.ErrUnresolvedFlow) NiceError {
	return NiceError{
		Code:    unableToResolveFlowCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Name": err.Error(),
		},
	}
}
func humanizeErrUndefinedReference(err references.ErrUndefinedReference) NiceError {
	details := map[string]interface{}{
		"Reference":  err.Reference,
		"Breakpoint": err.Breakpoint,
		"Path":       err.Path,
	}

	if err.Expression != nil {
		details["Expression"] = err.Expression.Position()
	}

	return NiceError{
		Code:    undefinedReferenceCode,
		Message: err.Error(),
		Details: details,
	}
}
func humanizeErrUndefinedResource(err references.ErrUndefinedResource) NiceError {
	var availableRefs []string
	for k, _ := range err.AvailableReferences {
		availableRefs = append(availableRefs, k)
	}

	return NiceError{
		Code:    undefinedResourceCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Reference":  err.Reference,
			"Breakpoint": err.Breakpoint,
			"KnownReferences": availableRefs,
		},
	}
}
func humanizeErrUnresolvedOnError(err references.ErrUnresolvedOnError) NiceError {
	return NiceError{
		Code:    unresolvedOnErrorCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"OnError": err.OnError,
		},
	}
}
func humanizeErrUnresolvedNode(err references.ErrUnresolvedNode) NiceError {
	return NiceError{
		Code:    unresolvedNodeCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Node": err.Node,
		},
	}
}
func humanizeErrUnresolvedParameterMap(err references.ErrUnresolvedParameterMap) NiceError {
	return NiceError{
		Code:    unresolvedParameterMapCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Parameter": err.Parameter,
		},
	}
}
func humanizeErrUnresolvedProperty(err references.ErrUnresolvedProperty) NiceError {
	return NiceError{
		Code:    unresolvedPropertyCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Property": err.Property,
		},
	}
}
func humanizeErrUnresolvedParams(err references.ErrUnresolvedParams) NiceError {
	return NiceError{
		Code:    unresolvedParamsCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Params": err.Params,
		},
	}
}
func humanizeErrUnresolvedCall(err references.ErrUnresolvedCall) NiceError {
	return NiceError{
		Code:    UnresolvedCallCode,
		Message: err.Error(),
		Details: map[string]interface{}{
			"Call": err.Call,
		},
	}
}
