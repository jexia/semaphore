package conditions

import (
	"regexp"

	"github.com/Knetic/govaluate"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/labels"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
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
		expressionTemplate, err := template.Parse(ctx, "", ref)
		if err != nil {
			return nil, err
		}

		// TODO Better use ParseTemplateProperty() when it is moved (see issue #194)
		param := &specs.Property{
			Name:     ref,
			Path:     "",
			Template: expressionTemplate,
		}

		if expressionTemplate.Reference == nil {
			param.Label = labels.Optional
		} else {
			// TODO Only Template.Reference got their Raw set (other types not), does this make sense?!
			param.Raw = ref
		}

		params.Params[ref] = param
	}

	result := &specs.Condition{
		RawExpression: raw,
		Evaluable:     expression,
		Params:        params,
	}

	return result, nil
}
