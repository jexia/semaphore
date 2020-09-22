package hcl

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// DefaultOnError sets the default values for not defined properties
func DefaultOnError(err *specs.OnError) {
	if err == nil {
		err = &specs.OnError{}
	}

	if err.Status == nil {
		err.Status = &specs.Property{
			Type:  types.Int64,
			Label: labels.Optional,
			Reference: &specs.PropertyReference{
				Resource: "error",
				Path:     "status",
			},
		}
	}

	if err.Message == nil {
		err.Message = &specs.Property{
			Type:  types.String,
			Label: labels.Optional,
			Reference: &specs.PropertyReference{
				Resource: "error",
				Path:     "message",
			},
		}
	}
}

// MergeOnError merges the right on error specs into the left on error
func MergeOnError(left *specs.OnError, right *specs.OnError) {
	if left == nil || right == nil {
		return
	}

	if left.Message == nil {
		left.Message = right.Message.Clone()
	}

	if left.Status == nil {
		left.Status = right.Status.Clone()
	}

	if len(left.Params) == 0 {
		left.Params = make(map[string]*specs.Property, len(right.Params))

		for key, param := range right.Params {
			left.Params[key] = param.Clone()
		}
	}

	if left.Response == nil {
		left.Response = right.Response.Clone()
	}
}

// ResolveErrors clones any previously defined error objects or error handles
func ResolveErrors(flows specs.FlowListInterface, err *specs.ParameterMap) {
	for _, flow := range flows {
		DefaultOnError(flow.GetOnError())

		if flow.GetOnError().Response == nil {
			flow.GetOnError().Response = err.Clone()
		}

		for _, node := range flow.GetNodes() {
			if node.OnError == nil {
				node.OnError = flow.GetOnError().Clone()
				continue
			}

			MergeOnError(node.OnError, flow.GetOnError())
		}
	}
}

// ErrMultiValueReource occurs when resource has more than one request and or flow
type ErrMultiValueReource struct {
	wrapErr
	Name string
}

// Error returns a description of the given error as a string
func (e ErrMultiValueReource) Error() string {
	return fmt.Sprintf("only one request and or flow could be defined inside a single resource '%s'", e.Name)
}

// Prettify returns the prettified version of the given error
func (e ErrMultiValueReource) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "MultiValueReource",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Name": e.Name,
		},
	}
}

// ErrPathNotFound occurs when resource has more than one request and or flow
type ErrPathNotFound struct {
	wrapErr
	Path string
}

// Error returns a description of the given error as a string
func (e ErrPathNotFound) Error() string {
	return fmt.Sprintf("unable to resolve path, no files found '%s'", e.Path)
}

// Prettify returns the prettified version of the given error
func (e ErrPathNotFound) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "PathNotFound",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path": e.Path,
		},
	}
}
