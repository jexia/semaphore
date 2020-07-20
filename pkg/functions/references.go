package functions

import (
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/references"
	"github.com/jexia/maestro/pkg/specs"
)

// DefineFunctions defined all properties within the given functions
func DefineFunctions(ctx instance.Context, functions Stack, node *specs.Node, flow specs.FlowResourceManager) error {
	if functions == nil {
		return nil
	}

	for _, function := range functions {
		if function.Arguments != nil {
			for _, arg := range function.Arguments {
				references.DefineProperty(ctx, node, arg, flow)
			}
		}

		if function.Returns == nil {
			continue
		}

		references.DefineProperty(ctx, node, function.Returns, flow)
	}

	return nil
}
