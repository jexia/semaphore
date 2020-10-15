package functions

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// DefineFunctions defined all properties within the given functions
func DefineFunctions(ctx *broker.Context, functions Stack, node *specs.Node, flow specs.FlowInterface) error {
	if functions == nil {
		return nil
	}

	for _, function := range functions {
		if function.Arguments != nil {
			for _, arg := range function.Arguments {
				if err := references.ResolveProperty(ctx, node, arg, flow); err != nil {
					return err
				}
			}
		}

		if function.Returns == nil {
			continue
		}

		if err := references.ResolveProperty(ctx, node, function.Returns, flow); err != nil {
			return err
		}
	}

	return nil
}
