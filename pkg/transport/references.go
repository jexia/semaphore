package transport

import (
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/dependencies"
	"github.com/jexia/maestro/pkg/specs/references"
)

// DefineCaller defineds the types for the given transport caller
func DefineCaller(ctx instance.Context, node *specs.Node, manifest *specs.FlowsManifest, call Call, manager specs.FlowResourceManager) (err error) {
	ctx.Logger(logger.Core).Info("Defining caller references")

	method := call.GetMethod(node.Call.Method)
	for _, prop := range method.References() {
		err = references.DefineProperty(ctx, node, prop, manager)
		if err != nil {
			return err
		}

		dependencies.ResolvePropertyReferences(prop, node.DependsOn)
		err = dependencies.ResolveNode(manager, node, make(map[string]*specs.Node))
		if err != nil {
			return err
		}
	}

	return nil
}
