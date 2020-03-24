package metadata

import (
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// NewManager constructs a new metadata manager for the given resource
func NewManager(resource string, params *specs.ParameterMap) *Manager {
	return &Manager{
		Resource: specs.JoinPath(resource, specs.ResourceHeader),
		Params:   params,
	}
}

// Manager represents a metadata manager for a given resource
type Manager struct {
	Resource string
	Params   *specs.ParameterMap
}

// Marshal attempts to marshal the given metadata specs from the given refs store
func (manager *Manager) Marshal(store *refs.Store) MD {
	if manager == nil || manager.Params == nil {
		return make(MD, 0)
	}

	result := make(MD, len(manager.Params.Header))
	for key, property := range manager.Params.Header {
		value := property.Default

		if property.Reference != nil {
			ref := store.Load(property.Reference.Resource, property.Reference.Path)
			if ref != nil {
				value = ref.Value
			}
		}

		if value == nil {
			continue
		}

		result[key] = value.(string)
	}

	return result
}

// Unmarshal unmarshals the given transport metadata into the given reference store
func (manager *Manager) Unmarshal(metadata MD, store *refs.Store) {
	for key, value := range metadata {
		ref := refs.New(key)
		ref.Value = value
		store.StoreReference(manager.Resource, ref)
	}
}
