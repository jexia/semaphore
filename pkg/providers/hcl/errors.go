package hcl

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

var errNonScalarType = errors.New("non scalar type")

type errUnknownPopertyType string

func (e errUnknownPopertyType) Error() string {
	return fmt.Sprintf("unsupported property type %q", string(e))
}

// DefaultOnError sets the default values for not defined properties
func DefaultOnError(err *specs.OnError) {
	if err == nil {
		err = &specs.OnError{}
	}

	if err.Status == nil {
		err.Status = &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Reference: &specs.PropertyReference{
					Resource: "error",
					Path:     "status",
				},
				Scalar: &specs.Scalar{
					Type: types.Int64,
				},
			},
		}
	}

	if err.Message == nil {
		err.Message = &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Reference: &specs.PropertyReference{
					Resource: "error",
					Path:     "message",
				},
				Scalar: &specs.Scalar{
					Type: types.String,
				},
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
