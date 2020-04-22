package checks

import (
	"sync"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/trace"
)

// ManifestDuplicates checks for duplicate definitions
func ManifestDuplicates(ctx instance.Context, manifest *specs.FlowsManifest) error {
	ctx.Logger(logger.Core).Info("Checking manifest duplicates")

	tracker := sync.Map{}

	for _, flow := range manifest.Flows {
		_, duplicate := tracker.LoadOrStore(flow.Name, flow)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate flow '%s'", flow.Name))
		}

		err := NodeDuplicates(ctx, flow.Name, flow.Nodes)
		if err != nil {
			return err
		}
	}

	for _, proxy := range manifest.Proxy {
		_, duplicate := tracker.LoadOrStore(proxy.Name, proxy)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate flow '%s'", proxy.Name))
		}

		err := NodeDuplicates(ctx, proxy.Name, proxy.Nodes)
		if err != nil {
			return err
		}
	}

	return nil
}

// NodeDuplicates checks for duplicate definitions
func NodeDuplicates(ctx instance.Context, flow string, nodes []*specs.Node) error {
	ctx.Logger(logger.Core).Info("Checking flow duplicates")

	calls := sync.Map{}

	for _, node := range nodes {
		_, duplicate := calls.LoadOrStore(node.Name, node)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate resource '%s' in flow '%s'", node.Name, flow))
		}
	}

	return nil
}
