package metadata

import (
	"strings"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"go.uber.org/zap"
)

// NewManager constructs a new metadata manager for the given resource.
func NewManager(ctx *broker.Context, resource string, params specs.Header) *Manager {
	return &Manager{
		Context:  ctx,
		Resource: template.JoinPath(resource, template.HeaderResource),
		Params:   params,
	}
}

// Manager represents a metadata manager for a given resource
type Manager struct {
	Context  *broker.Context
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

		logger.Debug(manager.Context, "Marshalling header property", zap.String("key", key))

		if value == nil {
			logger.Debug(manager.Context, "Header property is empty", zap.String("key", key))
			continue
		}

		result[key] = value.(string)
	}

	return result
}

// Unmarshal unmarshals the given transport metadata into the given reference store
func (manager *Manager) Unmarshal(metadata MD, store references.Store) {
	for key, value := range metadata {
		_, has := manager.Params[key]
		if !has {
			continue
		}

		ref := &references.Reference{
			Path:  strings.ToLower(key),
			Value: value,
		}

		logger.Debug(manager.Context, "Unmarshalling header property", zap.String("key", key))
		store.StoreReference(manager.Resource, ref)
	}
}
