package flow

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/refs"
	"github.com/jexia/semaphore/pkg/specs"
)

// NewCondition constructs a new condition of the given functions stack and specs condition
func NewCondition(stack functions.Stack, spec *specs.Condition) *Condition {
	return &Condition{
		stack:     stack,
		condition: spec,
	}
}

// Condition represents a condition which could be evaluated and results in a boolean
type Condition struct {
	stack     functions.Stack
	condition *specs.Condition
}

// Eval evaluates the given condition with the given reference store
func (condition *Condition) Eval(ctx instance.Context, store refs.Store) (bool, error) {
	err := ExecuteFunctions(condition.stack, store)
	if err != nil {
		return false, err
	}

	parameters := make(map[string]interface{}, len(condition.condition.Params.Params))
	for key, param := range condition.condition.Params.Params {
		value := param.Default

		if param.Reference != nil {
			ref := store.Load(param.Reference.Resource, param.Reference.Path)
			if ref != nil {
				value = ref.Value
			}
		}

		parameters[key] = value
	}

	ctx.Logger(logger.Flow).WithField("parameters", parameters).Debug("Evaluating comparison")
	result, err := condition.condition.Expression.Evaluate(parameters)
	if err != nil {
		return false, err
	}

	pass, is := result.(bool)
	if !is {
		return true, nil
	}

	return pass, nil
}
