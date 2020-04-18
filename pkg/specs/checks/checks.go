package checks

import (
	"sync"

	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/trace"
)

// ManifestDuplicates checks for duplicate definitions
func ManifestDuplicates(ctx instance.Context, manifest *specs.FlowsManifest) error {
	ctx.Logger(logger.Core).Info("Checking manifest duplicates")

	flows := sync.Map{}

	for _, flow := range manifest.Flows {
		_, duplicate := flows.LoadOrStore(flow.Name, flow)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate flow '%s'", flow.Name))
		}

		err := FlowDuplicates(ctx, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// FlowDuplicates checks for duplicate definitions
func FlowDuplicates(ctx instance.Context, flow *specs.Flow) error {
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
