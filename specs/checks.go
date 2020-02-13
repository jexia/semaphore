package specs

import (
	"sync"

	"github.com/jexia/maestro/specs/trace"
)

// CheckManifestDuplicates checks for duplicate definitions
func CheckManifestDuplicates(manifest *Manifest) error {
	callers := sync.Map{}
	endpoints := sync.Map{}
	flows := sync.Map{}
	services := sync.Map{}

	for _, caller := range manifest.Callers {
		_, duplicate := callers.LoadOrStore(caller.Name, caller)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate caller '%s'", manifest.File.Path, caller.Name))
		}
	}

	for _, endpoint := range manifest.Endpoints {
		_, duplicate := endpoints.LoadOrStore(endpoint.Flow, endpoint)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate flow endpoint '%s'", manifest.File.Path, endpoint.Flow))
		}
	}

	for _, service := range manifest.Services {
		_, duplicate := services.LoadOrStore(service.Alias, service)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate service alias '%s'", manifest.File.Path, service.Alias))
		}
	}

	for _, flow := range manifest.Flows {
		_, duplicate := flows.LoadOrStore(flow.Name, flow)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate flow '%s'", manifest.File.Path, flow.Name))
		}

		err := CheckFlowDuplicates(manifest, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckFlowDuplicates checks for duplicate definitions
func CheckFlowDuplicates(manifest *Manifest, flow *Flow) error {
	calls := sync.Map{}

	for _, call := range flow.Calls {
		_, duplicate := calls.LoadOrStore(call.Name, call)
		if duplicate {
			return trace.New(trace.WithMessage("%s duplicate call '%s' in flow '%s'", manifest.File.Path, call.Name, flow.Name))
		}
	}

	return nil
}
