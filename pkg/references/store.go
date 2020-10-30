package references

import (
	"fmt"
	"sync"

	"github.com/jexia/semaphore/pkg/lookup"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// Store is a key/value store capable of holding reference values.
// Reference values could be fetched by providing the absolute path of a property.
// The reference store is used inside flows and codecs to track properties and their values.
// Property keys are delimited with a simple dot-notation (ex: meta.key).
//
// Arrays are stored by defining the index and property path (ex: items[0].key).
// When defining an array or object make sure to define the length.
// The length of objects and properties are used inside implementations as a reference.
type Store interface {
	// Store stores the reference using the given path as key
	Store(path string, reference *Reference)
	// Load attempts to load the reference for the given path.
	// If no reference has been found is a nil value returned.
	Load(path string) *Reference
	// Define defines the length of a array or object at the given path.
	// Any previously defined lengths for the given path will be overridden.
	Define(path string, length int)
	// Length returns the length of the given object or array at the given path
	Length(path string) int
}

// Reference represents a value reference
type Reference struct {
	Value interface{}
	Enum  *int32
}

// NewStore constructs a new store and allocates the references for the given length
func NewStore(size int) Store {
	return &store{
		lengths: make(map[string]int, size),
		values:  make(map[string]*Reference, size),
	}
}

type store struct {
	lengths map[string]int
	values  map[string]*Reference
	mutex   sync.RWMutex
}

// Store stores the given value inside the references store
func (store *store) Store(path string, reference *Reference) {
	store.mutex.Lock()
	store.values[path] = reference
	store.mutex.Unlock()
}

// Load attempts to load the defined value for the given path.
// If the path ends with a self reference (.) is it trimmed off.
func (store *store) Load(path string) *Reference {
	if len(path) > 0 && string(path[len(path)-1]) == lookup.SelfRef {
		path = path[:len(path)-1]
	}

	store.mutex.RLock()
	ref, has := store.values[path]
	store.mutex.RUnlock()
	if !has {
		return nil
	}

	return ref
}

// Define attempts to load the defined value for the given path
func (store *store) Define(path string, length int) {
	store.mutex.Lock()
	store.lengths[path] = length
	store.mutex.Unlock()
}

// Load attempts to load the defined value for the given path
func (store *store) Length(path string) int {
	store.mutex.RLock()
	length, has := store.lengths[path]
	store.mutex.RUnlock()
	if !has {
		return 0
	}

	return length // returns 0 if not defined
}

// Tracker tracks the index positions of arrays.
// Paths represent arrays or objects stored inside a reference store.
// Paths defined inside references could be resolved to include the current items index.
// These indexes are required since paths inside the reference store are absolute.
// Trackers could be nested to track multiple arrays on different or nested paths.
//
// The tracker does not perform any mutex locking.
// Most usecases/implementations are not concurrent.
//
// example:
// track := tracker.Track("items", 0) // sets the current index of "items" at 0
// path := track.Resolve("items.key") // returns: "items[0].key"
// track.Next("items")                // increases the index of "items" by one
type Tracker interface {
	// Track includes the given path and index to be tracked.
	// Trackers could be nested to track multiple or nested arrays on different indexes.
	Track(path string, index int)
	// Resolve resolves the given path to include all tracked array indexes found inside the given path.
	// ex: "items.key" - "items[0].key"
	Resolve(path string) string
	// Next increases the index of the counter at the given path
	Next(path string) int
	// NOTE: remove tracked paths?
}

// NewTracker constructs a new position tracker
func NewTracker() Tracker {
	return &tracker{
		positions: map[string]int{},
	}
}

type tracker struct {
	positions map[string]int
}

func (t *tracker) Track(path string, index int) {
	t.positions[path] = index
}

func (t *tracker) Resolve(path string) string {
	if len(t.positions) == 0 {
		return path
	}

	// TODO: replace this implementation with a Radix tree
	result := ""
	lookup := ""

	parts := template.SplitPath(path)

	for _, part := range parts {
		result = template.JoinPath(result, part)
		lookup = template.JoinPath(lookup, part)

		index, has := t.positions[lookup]
		if !has {
			continue
		}

		result += fmt.Sprintf("[%d]", index)
	}

	return result
}

func (t *tracker) Next(path string) int {
	t.positions[path]++
	return t.positions[path]
}

// Index returns the path and the provided index as a path
func Index(path string, index int) string {
	return fmt.Sprintf("%s[%d]", path, index)
}

// StoreValues stores the given values to the reference store
func StoreValues(store Store, tracker Tracker, path string, values map[string]interface{}) {
	store.Define(path, len(values))

	for key, value := range values {
		path := template.JoinPath(path, key)
		keys, is := value.(map[string]interface{})
		if is {
			StoreValues(store, tracker, path, keys)
			continue
		}

		repeated, is := value.([]map[string]interface{})
		if is {
			NewRepeatingMessages(store, tracker, path, repeated)
			continue
		}

		values, is := value.([]interface{})
		if is {
			NewRepeating(store, tracker, path, values)
			continue
		}

		enum, is := value.(*EnumVal)
		if is {
			store.Store(tracker.Resolve(path), &Reference{
				Enum: &enum.pos,
			})
			continue
		}

		store.Store(tracker.Resolve(path), &Reference{
			Value: value,
		})
	}
}

// NewRepeatingMessages appends the given repeating messages to the given reference
func NewRepeatingMessages(store Store, tracker Tracker, path string, values []map[string]interface{}) {
	tracker.Track(path, 0)

	for _, values := range values {
		StoreValues(store, tracker, path, values)
		tracker.Next(path)
	}

	store.Define(path, len(values))
}

// NewRepeating appends the given repeating values to the given reference
func NewRepeating(store Store, tracker Tracker, path string, values []interface{}) {
	store.Define(path, len(values))
	tracker.Track(path, 0)

	for _, value := range values {
		enum, is := value.(*EnumVal)
		if is {
			store.Store(tracker.Resolve(path), &Reference{
				Enum: &enum.pos,
			})

			tracker.Next(path)
			continue
		}

		store.Store(tracker.Resolve(path), &Reference{
			Value: value,
		})

		tracker.Next(path)
	}
}

// NewPrefixStore fixes all writes and reads from the given store on the set resource and prefix path
func NewPrefixStore(store Store, prefix string) Store {
	return &prefixStore{
		path:  prefix,
		store: store,
	}
}

// prefixStore creates a sandbox where all resources stored are forced into the set resource and prefix
type prefixStore struct {
	path  string
	store Store
}

// Load attempts to load the defined value for the given resource and path
func (prefix *prefixStore) Load(path string) *Reference {
	return prefix.store.Load(path)
}

// Store stores the given value inside the references store
func (prefix *prefixStore) Store(path string, reference *Reference) {
	prefix.store.Store(template.JoinPath(prefix.path, path), reference)
}

// Define attempts to load the defined value for the given path.
func (prefix *prefixStore) Define(path string, length int) {
	prefix.store.Define(template.JoinPath(prefix.path, path), length)
}

// Load attempts to load the defined value for the given resource and path
func (prefix *prefixStore) Length(path string) int {
	return prefix.store.Length(path)
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
