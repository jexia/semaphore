package refs

import (
	"sync"

	"github.com/jexia/maestro/specs"
)

// New constructs a new reference with the given path
func New(path string) *Reference {
	return &Reference{
		Path: path,
	}
}

// Reference represents a value reference
type Reference struct {
	Path     string
	Value    interface{}
	Repeated []*Store
	mutex    sync.Mutex
}

// Repeating prepares the given reference to store repeating values
func (reference *Reference) Repeating(size int) {
	reference.Repeated = make([]*Store, size)
}

// Set sets the given repeating value reference on the given index
func (reference *Reference) Set(index int, val *Store) {
	reference.mutex.Lock()
	reference.Repeated[index] = val
	reference.mutex.Unlock()
}

// NewStore constructs a new store and allocates the references for the given length
func NewStore(size int) *Store {
	return &Store{
		values: make(map[string]*Reference, size),
	}
}

// Store references
type Store struct {
	values map[string]*Reference
	mutex  sync.Mutex
}

// StoreReference stores the given resource, path and value inside the references store
func (store *Store) StoreReference(resource string, reference *Reference) {
	hash := resource + reference.Path
	store.mutex.Lock()
	store.values[hash] = reference
	store.mutex.Unlock()
}

// Load attempts to load the defined value for the given resource and path
func (store *Store) Load(resource string, path string) *Reference {
	hash := resource + path
	store.mutex.Lock()
	ref, has := store.values[hash]
	store.mutex.Unlock()
	if !has {
		return nil
	}

	return ref
}

// StoreValues stores the given values to the reference store
func (store *Store) StoreValues(resource string, path string, values map[string]interface{}) {
	for key, val := range values {
		path := specs.JoinPath(path, key)
		values, is := val.(map[string]interface{})
		if is {
			store.StoreValues(resource, path, values)
			continue
		}

		repeated, is := val.([]interface{})
		if is {
			reference := New(path)
			store.NewRepeating(resource, path, reference, repeated)
			store.StoreReference(resource, reference)
			continue
		}

		store.StoreValue(resource, path, val)
	}
}

// StoreValue stores the given value for the given resource and path
func (store *Store) StoreValue(resource string, path string, value interface{}) {
	reference := New(path)
	reference.Value = value

	store.StoreReference(resource, reference)
}

// NewRepeating appends the given repeating values to the given reference
func (store *Store) NewRepeating(resource string, path string, reference *Reference, values []interface{}) {
	reference.Repeating(len(values))

	for index, value := range values {
		values, is := value.(map[string]interface{})
		if !is {
			continue
		}

		store := &Store{
			values: make(map[string]*Reference, len(values)),
		}

		store.StoreValues(resource, path, values)
		reference.Set(index, store)
	}
}
