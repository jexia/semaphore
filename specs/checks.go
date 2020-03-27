package specs

import (
	"sync"

	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs/trace"
)

// CheckManifestDuplicates checks for duplicate definitions
func CheckManifestDuplicates(ctx instance.Context, manifest *Manifest) error {
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
func CheckFlowDuplicates(ctx instance.Context, flow *Flow) error {
	ctx.Logger(logger.Core).Info("Checking flow duplicates")

	calls := sync.Map{}

	for _, call := range flow.Nodes {
		_, duplicate := calls.LoadOrStore(call.Name, call)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate call '%s' in flow '%s'", call.Name, flow.Name))
		}
	}

	return nil
}
