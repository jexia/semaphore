package specs

import (
	"sync"

	"github.com/jexia/maestro/specs/trace"
	log "github.com/sirupsen/logrus"
)

// CheckManifestDuplicates checks for duplicate definitions
func CheckManifestDuplicates(manifest *Manifest) error {
	log.Info("Checking manifest duplicates")

	flows := sync.Map{}

	for _, flow := range manifest.Flows {
		_, duplicate := flows.LoadOrStore(flow.Name, flow)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate flow '%s'", flow.Name))
		}

		err := CheckFlowDuplicates(flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckFlowDuplicates checks for duplicate definitions
func CheckFlowDuplicates(flow *Flow) error {
	log.Info("Checking flow duplicates")

	calls := sync.Map{}

	for _, call := range flow.Nodes {
		_, duplicate := calls.LoadOrStore(call.Name, call)
		if duplicate {
			return trace.New(trace.WithMessage("duplicate call '%s' in flow '%s'", call.Name, flow.Name))
		}
	}

	return nil
}
