package conditions

import (
	"regexp"

	"github.com/Knetic/govaluate"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

var templatePattern = regexp.MustCompile("{{ *([^}}]+?) *}}")

// NewEvaluableExpression constructs a new condition out of the given expression
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
		prop, err := template.Parse(ctx, "", ref, ref)
		if err != nil {
			return nil, err
		}

		params.Params[ref] = prop
	}

	result := &specs.Condition{
		RawExpression: raw,
		Evaluable:     expression,
		Params:        params,
	}

	return result, nil
}
