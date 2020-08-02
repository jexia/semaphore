package metadata

import (
	"strings"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// NewManager constructs a new metadata manager for the given resource.
func NewManager(ctx instance.Context, resource string, params specs.Header) *Manager {
	return &Manager{
		Context:  ctx,
		Resource: template.JoinPath(resource, template.HeaderResource),
		Params:   params,
	}
}

// Manager represents a metadata manager for a given resource
type Manager struct {
	Context  instance.Context
	Resource string
	Params   specs.Header
}

// Marshal attempts to marshal the given metadata specs from the given refs store
func (manager *Manager) Marshal(store references.Store) MD {
	if manager == nil || manager.Params == nil {
		return make(MD, 0)
	}

	result := make(MD, len(manager.Params))
	for key, property := range manager.Params {
		value := property.Default

		if property.Reference != nil {
			ref := store.Load(property.Reference.Resource, property.Reference.Path)
			if ref != nil {
				value = ref.Value
			}
		}

		manager.Context.Logger(logger.Flow).WithField("key", key).Debug("Marshalling header property")

		if value == nil {
			manager.Context.Logger(logger.Flow).WithField("key", key).Debug("Header property is empty")
			continue
		}

		result[key] = value.(string)
	}

	return result
}

// Unmarshal unmarshals the given transport metadata into the given reference store
func (manager *Manager) Unmarshal(metadata MD, store references.Store) {
	for key, value := range metadata {
		ref := references.NewReference(strings.ToLower(key))
		ref.Value = value

		manager.Context.Logger(logger.Flow).WithField("key", key).Debug("Unmarshalling header property")

		store.StoreReference(manager.Resource, ref)
	}
}
