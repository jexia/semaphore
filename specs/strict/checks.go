package strict

import (
	"sync"

	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
)

// CheckManifestDuplicates checks for duplicate definitions
func CheckManifestDuplicates(ctx instance.Context, manifest *specs.Manifest) error {
	ctx.Logger(logger.Core).Info("Checking manifest duplicates")

	flows := sync.Map{}

	for _, flow := range manifest.Flows {
		_, duplicate := flows.LoadOrStore(flow.Name, flow)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate flow '%s'", flow.Name))
		}

		err := CheckFlowDuplicates(ctx, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckFlowDuplicates checks for duplicate definitions
func CheckFlowDuplicates(ctx instance.Context, flow *specs.Flow) error {
	ctx.Logger(logger.Core).Info("Checking flow duplicates")

	calls := sync.Map{}

	for _, node := range flow.Nodes {
		_, duplicate := calls.LoadOrStore(node.Name, node)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate resource '%s' in flow '%s'", node.Name, flow.Name))
		}
	}

	return nil
}
