package formencoded

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Name represents the codec
const Name = "form-urlencoded"

// NewConstructor constructs a new www-form-urlencoded constructor
func NewConstructor() *Constructor {
	return &Constructor{}
}

// Constructor is capable of constructing new codec managers for the given resource and specs
type Constructor struct {
}

// Name returns the name of the www-form-urlencoded codec constructor
func (constructor *Constructor) Name() string {
	return Name
}

// New constructs a new www-form-urlencoded codec manager
func (constructor *Constructor) New(resource string, specs *specs.ParameterMap) (codec.Manager, error) {
	if specs == nil {
		return nil, trace.New(trace.WithMessage("no object specs defined"))
	}

	return &Manager{
		resource: resource,
		specs:    specs.Property,
	}, nil
}

// Manager manages a specs object and allows to encode/decode messages
type Manager struct {
	resource string
	specs    *specs.Property
}

// Name returns the proto codec name
func (manager *Manager) Name() string {
	return Name
}

// Property returns the manager property which is used to marshal and unmarshal data
func (manager *Manager) Property() *specs.Property {
	return manager.specs
}

// Marshal marshals the given reference store into a www-form-urlencoded message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store
func (manager *Manager) Marshal(refs references.Store) (io.Reader, error) {
	if manager.specs == nil {
		return bytes.NewReader([]byte{}), nil
	}

	encoder := url.Values{}
	encode(encoder, "", manager.specs.Name, refs, manager.specs)

	bb := []byte(encoder.Encode())
	return bytes.NewReader(bb), nil
}

func encode(encoder url.Values, root string, name string, refs references.Store, prop *specs.Property) {
	path := template.JoinPath(root, name)

	switch {
	case prop.Message != nil:
		for key, nested := range prop.Message.Properties {
			encode(encoder, path, key, refs, nested)
		}
	case prop.Repeated != nil:
		if prop.Reference == nil {
			return
		}

		ref := refs.Load(prop.Reference.Resource, prop.Reference.Path)
		if ref == nil {
			break
		}

		for index, repeated := range ref.Repeated {
			current := fmt.Sprintf("%s[%d]", path, index)
			encode(encoder, current, "", repeated, prop.Repeated.Property)
		}
	case prop.Enum != nil:
		ref := refs.Load(prop.Reference.Resource, prop.Reference.Path)
		if ref == nil {
			break
		}

		key := prop.Enum.Positions[*ref.Enum]
		AddTypeKey(encoder, path, types.Enum, key)
	case prop.Scalar != nil:
		val := prop.Scalar.Default

		ref := refs.Load(prop.Reference.Resource, prop.Reference.Path)
		if ref != nil {
			val = ref.Value
		}

		if val == nil {
			break
		}

		AddTypeKey(encoder, path, prop.Scalar.Type, val)
	}
}

// Unmarshal unmarshals the given www-form-urlencoded io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs references.Store) error {
	if manager.specs == nil {
		return nil
	}

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if len(bb) == 0 {
		return nil
	}

	values, err := url.ParseQuery(string(bb))
	if err != nil {
		return err
	}

	for key, values := range values {
		if err := decodeElement(manager.resource, 0, strings.Split(key, "."), values, manager.specs.Nested, refs); err != nil {
			return err
		}
	}

	return nil
}

func decodeElement(resource string, pos int, path []string, values []string, schema specs.PropertyList, refs references.Store) error {
	propName := path[pos]

	if schema == nil {
		return errNilSchema
	}

	prop := schema.Get(propName)
	if prop == nil {
		return errUndefinedProperty(propName)
	}

	ref := &references.Reference{
		Path: prop.Path,
	}

	switch prop.Label {
	case labels.Repeated:
		for _, raw := range values {
			store := references.NewReferenceStore(0)

			switch prop.Type {
			case types.Message:
				if len(path) > pos+1 {
					if err := decodeElement(resource, pos+1, path, []string{raw}, prop.Nested, store); err != nil {
						return err
					}
				}
			case types.Enum:
				enum := prop.Enum.Keys[raw]
				if enum != nil {
					store.StoreEnum("", "", enum.Position)
				}
			default:
				value, err := types.DecodeFromString(raw, prop.Type)
				if err != nil {
					return err
				}

				store.StoreValue("", "", value)
			}

			ref.Append(store)
		}
	case labels.Optional, labels.Required:
		switch prop.Type {
		case types.Message:
			if len(path) > pos+1 {
				return decodeElement(resource, pos+1, path, values, prop.Nested, refs)
			}
		case types.Enum:
			enum := prop.Enum.Keys[values[0]]
			if enum != nil {
				ref.Enum = &enum.Position
			}
		default:
			value, err := types.DecodeFromString(values[0], prop.Type)
			if err != nil {
				return err
			}

			ref.Value = value
		}
	default:
		return errUnknownLabel(prop.Label)
	}

	refs.StoreReference(resource, ref)

	return nil
}
