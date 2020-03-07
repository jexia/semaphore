package protocol

import (
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// A Header represents the key-value pairs.
type Header map[string]string

// Clone returns a copy of h or nil if h is nil.
func (h Header) Clone() Header {
	return h
}

// Del deletes the values associated with key.
func (h Header) Del(key string) {
	delete(h, key)
}

// Get gets the first value associated with the given key. If there are no values associated with the key, Get returns "".
func (h Header) Get(key string) string {
	return h[key]
}

// Set sets the header entries associated with key to the single element value. It replaces any existing values associated with key.
func (h Header) Set(key, value string) {
	h[key] = value
}

// NewHeaderManager constructs a new header manager for the given resource
func NewHeaderManager(resource string, params *specs.ParameterMap) *HeaderManager {
	return &HeaderManager{
		Resource: specs.JoinPath(resource, specs.ResourceHeader),
		Params:   params,
	}
}

// HeaderManager represents a header manager for a given resource
type HeaderManager struct {
	Resource string
	Params   *specs.ParameterMap
}

// Marshal attempts to marshal the given header specs from the given refs store
func (manager *HeaderManager) Marshal(store *refs.Store) Header {
	if manager.Params == nil {
		return make(Header, 0)
	}

	result := make(Header, len(manager.Params.Header))
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
func (manager *HeaderManager) Unmarshal(header Header, store *refs.Store) {
	for key, value := range header {
		ref := refs.New(key)
		ref.Value = value
		store.StoreReference(manager.Resource, ref)
	}
}
