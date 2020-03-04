package specs

import (
	"sync"

	"github.com/jexia/maestro/specs/trace"
	log "github.com/sirupsen/logrus"
)

// CheckManifestDuplicates checks for duplicate definitions
func CheckManifestDuplicates(file string, manifest *Manifest) error {
	log.Info("Checking manifest duplicates")

	endpoints := sync.Map{}
	flows := sync.Map{}
	services := sync.Map{}

	for _, endpoint := range manifest.Endpoints {
		_, duplicate := endpoints.LoadOrStore(endpoint.Flow, endpoint)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate flow endpoint '%s'", file, endpoint.Flow))
		}
	}

	for _, service := range manifest.Services {
		_, duplicate := services.LoadOrStore(service.Name, service)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate service '%s'", file, service.Name))
		}
	}

	for _, flow := range manifest.Flows {
		_, duplicate := flows.LoadOrStore(flow.Name, flow)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate flow '%s'", file, flow.Name))
		}

		err := CheckFlowDuplicates(file, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckFlowDuplicates checks for duplicate definitions
func CheckFlowDuplicates(file string, flow *Flow) error {
	log.Info("Checking flow duplicates")

	calls := sync.Map{}

	for _, call := range flow.Nodes {
		_, duplicate := calls.LoadOrStore(call.Name, call)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate call '%s' in flow '%s'", file, call.Name, flow.Name))
		}
	}

	return nil
}
