package hcl

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func TestExpressionPosition(t *testing.T) {
	t.Run("get position", func(t *testing.T) {
		var (
			expected = "file:10"

			expression = Expression{
				hcl.StaticExpr(
					cty.StringVal("prop"),
					hcl.Range{
						Filename: "file",
						Start: hcl.Pos{
							Line: 10,
						},
					},
				),
			}
		)

		if actual := expression.Position(); actual != expected {
			t.Errorf("unexpected position %q should be %q", actual, expected)
		}
	})
}
