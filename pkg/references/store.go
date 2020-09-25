package references

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// Store represents the reference store interface
type Store interface {
	// StoreReference stores the given resource, path and value inside the references store
	StoreReference(resource string, reference *Reference)
	// Load attempts to load the defined value for the given resource and path
	Load(resource string, path string) *Reference
	// StoreValues stores the given values to the reference store
	StoreValues(resource string, path string, values map[string]interface{})
	// StoreValue stores the given value for the given resource and path
	StoreValue(resource string, path string, value interface{})
	// StoreEnum stores the given enum on the given path
	StoreEnum(resource string, path string, enum int32)
}

// Reference represents a value reference
type Reference struct {
	Path     string
	Value    interface{}
	Repeated []Store
	Enum     *int32
	mutex    sync.Mutex
}

func (reference *Reference) String() string {
	switch {
	case reference.Value != nil:
		return fmt.Sprintf("%s:<%T(%v)>", reference.Path, reference.Value, reference.Value)
	case reference.Repeated != nil:
		return fmt.Sprintf("%s:<array(%s)>", reference.Path, reference.Repeated)
	case reference.Enum != nil:
		return fmt.Sprintf("%s:<enum(%d)>", reference.Path, *reference.Enum)
	default:
		return fmt.Sprintf("%s:<empty>", reference.Path)
	}
}

// Repeating prepares the given reference to store repeating values
func (reference *Reference) Repeating(size int) {
	reference.Repeated = make([]Store, size)
}

// Append appends the given store to the repeating value reference.
// This method uses append, it is advised to use Set & Repeating when the length of the repeated message is known.
func (reference *Reference) Append(val Store) {
	reference.mutex.Lock()
	reference.Repeated = append(reference.Repeated, val)
	reference.mutex.Unlock()
}

// Set sets the given repeating value reference on the given index
func (reference *Reference) Set(index int, val Store) {
	reference.mutex.Lock()
	reference.Repeated[index] = val
	reference.mutex.Unlock()
}

// NewReferenceStore constructs a new store and allocates the references for the given length
func NewReferenceStore(size int) Store {
	return &store{
		values: make(map[string]*Reference, size),
	}
}

type store struct {
	values map[string]*Reference
	mutex  sync.Mutex
}

func (store *store) String() string {
	var (
		separated bool
		builder   strings.Builder
	)

	for key, ref := range store.values {
		if separated {
			builder.WriteString(", ")
		} else {
			separated = true
		}

		builder.WriteString(fmt.Sprintf("%s:[%s]", key, ref))
	}

	return builder.String()
}

// StoreReference stores the given resource, path and value inside the references store
func (store *store) StoreReference(resource string, reference *Reference) {
	hash := resource + reference.Path
	store.mutex.Lock()
	store.values[hash] = reference
	store.mutex.Unlock()
}

// Load attempts to load the defined value for the given resource and path
func (store *store) Load(resource string, path string) *Reference {
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
func (store *store) StoreValues(resource string, path string, values map[string]interface{}) {
	for key, val := range values {
		path := template.JoinPath(path, key)
		keys, is := val.(map[string]interface{})
		if is {
			store.StoreValues(resource, path, keys)
			continue
		}

		repeated, is := val.([]map[string]interface{})
		if is {
			reference := &Reference{
				Path: path,
			}

			store.NewRepeatingMessages(resource, path, reference, repeated)
			store.StoreReference(resource, reference)
			continue
		}

		values, is := val.([]interface{})
		if is {
			reference := &Reference{
				Path: path,
			}

			store.NewRepeating(resource, path, reference, values)
			store.StoreReference(resource, reference)
			continue
		}

		enum, is := val.(*EnumVal)
		if is {
			store.StoreEnum(resource, path, enum.pos)
			continue
		}

		store.StoreValue(resource, path, val)
	}
}

// StoreValue stores the given value for the given resource and path
func (store *store) StoreValue(resource string, path string, value interface{}) {
	reference := &Reference{
		Path:  path,
		Value: value,
	}

	store.StoreReference(resource, reference)
}

// StoreEnum stores the given enum for the given resource and path
func (store *store) StoreEnum(resource string, path string, enum int32) {
	reference := &Reference{
		Path: path,
		Enum: &enum,
	}

	store.StoreReference(resource, reference)
}

// NewRepeatingMessages appends the given repeating messages to the given reference
func (store *store) NewRepeatingMessages(resource string, path string, reference *Reference, values []map[string]interface{}) {
	reference.Repeating(len(values))

	for index, values := range values {
		store := NewReferenceStore(len(values))
		store.StoreValues(resource, path, values)
		reference.Set(index, store)
	}
}

// NewRepeating appends the given repeating values to the given reference
func (store *store) NewRepeating(resource string, path string, reference *Reference, values []interface{}) {
	reference.Repeating(len(values))

	for index, value := range values {
		store := NewReferenceStore(1)

		enum, is := value.(*EnumVal)
		if is {
			store.StoreEnum("", "", enum.pos)
			reference.Set(index, store)
			continue
		}

		store.StoreValue("", "", value)
		reference.Set(index, store)
	}
}

// NewPrefixStore fixes all writes and reads from the given store on the set resource and prefix path
func NewPrefixStore(store Store, resource string, prefix string) Store {
	return &PrefixStore{
		resource: resource,
		path:     prefix,
		store:    store,
	}
}

// PrefixStore creates a sandbox where all resources stored are forced into the set resource and prefix
type PrefixStore struct {
	resource string
	path     string
	store    Store
}

// Load attempts to load the defined value for the given resource and path
func (prefix *PrefixStore) Load(resource string, path string) *Reference {
	return prefix.store.Load(resource, path)
}

// StoreReference stores the given resource, path and value inside the references store
func (prefix *PrefixStore) StoreReference(resource string, reference *Reference) {
	reference.Path = template.JoinPath(prefix.path, reference.Path)
	prefix.store.StoreReference(prefix.resource, reference)
}

// StoreValues stores the given values to the reference store
func (prefix *PrefixStore) StoreValues(resource string, path string, values map[string]interface{}) {
	prefix.store.StoreValues(prefix.resource, template.JoinPath(prefix.path, path), values)
}

// StoreValue stores the given value for the given resource and path
func (prefix *PrefixStore) StoreValue(resource string, path string, value interface{}) {
	prefix.store.StoreValue(prefix.resource, template.JoinPath(prefix.path, path), value)
}

// StoreEnum stores the given enum for the given resource and path
func (prefix *PrefixStore) StoreEnum(resource string, path string, enum int32) {
	prefix.store.StoreEnum(prefix.resource, template.JoinPath(prefix.path, path), enum)
}

// Collection represents a map of property references
type Collection map[string]*specs.PropertyReference

// MergeLeft merges the references into the given reference
func (references Collection) MergeLeft(incoming ...Collection) {
	for _, refs := range incoming {
		for key, val := range refs {
			references[key] = val
		}
	}
}

// ParameterReferences returns all the available references inside the given parameter map
func ParameterReferences(params *specs.ParameterMap) Collection {
	result := make(map[string]*specs.PropertyReference)

	if params == nil {
		return Collection{}
	}

	if params.Header != nil {
		for _, prop := range params.Header {
			if prop.Reference != nil {
				result[prop.Reference.String()] = prop.Reference
			}
		}
	}

	if params.Property != nil {
		for key, prop := range PropertyReferences(params.Property) {
			result[key] = prop
		}
	}

	return result
}

// PropertyReferences returns the available references within the given property
func PropertyReferences(property *specs.Property) Collection {
	result := make(map[string]*specs.PropertyReference)

	if property.Reference != nil {
		result[property.Reference.String()] = property.Reference
	}

	switch {
	case property.Message != nil:
		for _, nested := range property.Message {
			for key, ref := range PropertyReferences(nested) {
				result[key] = ref
			}
		}

		break
	case property.Repeated != nil:
		for _, repeated := range property.Repeated {
			property := &specs.Property{
				Template: repeated,
			}

			for key, ref := range PropertyReferences(property) {
				result[key] = ref
			}
		}

		break
	}

	return result
}
