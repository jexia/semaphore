package checks

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrFlowDuplicate occurs when FlowDuplicates finds several flows with the same name
type ErrFlowDuplicate struct {
	Flow string
}

func (e ErrFlowDuplicate) Error() string {
	return fmt.Sprintf("duplicate flow '%s'", e.Flow)
}

func (e ErrFlowDuplicate) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Message: e.Error(),
		Details: map[string]interface{}{
			"flow": e.Flow,
		},
	}
}

// ErrResourceDuplicate occurs when NodeDuplicates finds several resources with the same name
type ErrResourceDuplicate struct {
	Flow     string
	Resource string
}

func (e ErrResourceDuplicate) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Original: nil,
		Message:  e.Error(),
		Details: map[string]interface{}{
			"flow":     e.Flow,
			"resource": e.Resource,
		},
	}
}

func (e ErrResourceDuplicate) Error() string {
	return fmt.Sprintf("duplicate resource '%s' in flow '%s'", e.Resource, e.Flow)
}

// ErrReservedKeyword occurs when a flow's name conflicts with a reserved keyword.
type ErrReservedKeyword struct {
	Flow string
}

func (e ErrReservedKeyword) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Message: e.Error(),
		Details: map[string]interface{}{
			"flow": e.Flow,
		},
	}
}

func (e ErrReservedKeyword) Error() string {
	return fmt.Sprintf("flow with the name '%s' is a reserved keyword", e.Flow)
}
