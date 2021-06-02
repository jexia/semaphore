package conditions

import (
	"regexp"

	"github.com/Knetic/govaluate"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

var templatePattern = regexp.MustCompile("{{ *([^}}]+?) *}}")

// NewEvaluableExpression constructs a new condition out of the given expression.
func NewEvaluableExpression(ctx *broker.Context, raw string) (*specs.Condition, error) {
	raw = templatePattern.ReplaceAllString(raw, "[$1]") // replace template tags with govaluate parameter escape characters

	expression, err := govaluate.NewEvaluableExpression(raw)
	if err != nil {
		return nil, err
	}

	params := &specs.ParameterMap{
		Params: map[string]*specs.Property{},
	}

	for _, ref := range expression.Vars() {
		expressionTemplate, err := specs.ParseTemplate(ctx, "", ref)
		if err != nil {
			return nil, err
		}

		params.Params[ref] = specs.ParseTemplateProperty("", ref, ref, expressionTemplate)
	}

	result := &specs.Condition{
		RawExpression: raw,
		Evaluable:     expression,
		Params:        params,
	}

	return result, nil
}
