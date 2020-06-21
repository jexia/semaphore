package flow

import (
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

func NewCondition(functions functions.Stack, spec *specs.Condition) *Condition {
	return &Condition{
		functions: functions,
		condition: spec,
	}
}

type Condition struct {
	functions functions.Stack
	condition *specs.Condition
}

func (condition *Condition) Eval(ctx instance.Context, store refs.Store) (bool, error) {
	parameters := make(map[string]interface{}, len(condition.condition.Params.Params))
	for key, param := range condition.condition.Params.Params {
		ref := store.Load(param.Reference.Resource, param.Reference.Path)
		if ref == nil {
			parameters[key] = nil
			continue
		}

		parameters[key] = ref.Value
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
