package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// Expression is a wrapper over hcl.Expression providing extra functionality.
type Expression struct{ hcl.Expression }

// Position returns the position for tracer.
func (expr Expression) Position() string {
	return fmt.Sprintf("%s:%d", expr.Range().Filename, expr.Range().Start.Line)
}
