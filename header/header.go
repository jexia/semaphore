package header

import (
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// Store represents the key-value pairs.
type Store map[string]string

// Clone returns a copy of h or nil if h is nil.
func (h Store) Clone() Store {
	return h
}

// Del deletes the values associated with key.
func (h Store) Del(key string) {
	delete(h, key)
}

// Get gets the first value associated with the given key. If there are no values associated with the key, Get returns "".
func (h Store) Get(key string) string {
	return h[key]
}

// Set sets the header entries associated with key to the single element value. It replaces any existing values associated with key.
func (h Store) Set(key, value string) {
	h[key] = value
}

// NewManager constructs a new header manager for the given resource
func NewManager(resource string, params *specs.ParameterMap) *Manager {
	return &Manager{
		Resource: specs.JoinPath(resource, specs.ResourceHeader),
		Params:   params,
	}
}

// Manager represents a header manager for a given resource
type Manager struct {
	Resource string
	Params   *specs.ParameterMap
}

// Marshal attempts to marshal the given header specs from the given refs store
func (manager *Manager) Marshal(store *refs.Store) Store {
	if manager == nil || manager.Params == nil {
		return make(Store, 0)
	}

	result := make(Store, len(manager.Params.Header))
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

// Unmarshal unmarshals the given protocol header into the given reference store
func (manager *Manager) Unmarshal(header Store, store *refs.Store) {
	for key, value := range header {
		ref := refs.New(key)
		ref.Value = value
		store.StoreReference(manager.Resource, ref)
	}
}
