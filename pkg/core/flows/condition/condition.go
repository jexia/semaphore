package condition

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
)

// New constructs a new flow condition of the given specs
func New(ctx instance.Context, mem functions.Collection, condition *specs.Condition) *flow.Condition {
	if condition == nil {
		return nil
	}

	stack := mem[condition.Params]
	return flow.NewCondition(stack, condition)
}
