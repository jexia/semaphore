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
		property: specs.Property,
	}, nil
}

// Manager manages a specs object and allows to encode/decode messages
type Manager struct {
	resource string
	property *specs.Property
}

// Name returns the proto codec name
func (manager *Manager) Name() string {
	return Name
}

// Property returns the manager property which is used to marshal and unmarshal data
func (manager *Manager) Property() *specs.Property {
	return manager.property
}

// Marshal marshals the given reference store into a www-form-urlencoded message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store
func (manager *Manager) Marshal(refs references.Store) (io.Reader, error) {
	if manager.property == nil {
		return bytes.NewReader([]byte{}), nil
	}

	encoder := url.Values{}
	path := template.JoinPath("", manager.property.Name)
	err := encode(encoder, path, refs, manager.property.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to encode %s: %w", manager.property.Name, err)
	}

	bb := []byte(encoder.Encode())
	return bytes.NewReader(bb), nil
}

// encode the template recursively.
//
// - encoded is passed by the reference and the function modifies the given encoded argument by adding new key-value pairs.
// - path is the current value path in the encoded results. Example: "user.name", "id", "users[0]"
// - store is the references store
// - tpl is the encoding template
//
// The producing key-value pair examples:
// user.name=bob&user.age=30&id=100
// users[0]=bob&users[1]=alice
func encode(encoded url.Values, path string, store references.Store, tpl specs.Template) error {
	var (
		ref *references.Reference
	)

	if tpl.Reference != nil {
		ref = store.Load(tpl.Reference.Resource, tpl.Reference.Path)
	}

	switch {
	case tpl.Message != nil:
		for fieldName, field := range tpl.Message {
			path := template.JoinPath(path, fieldName)
			err := encode(encoded, path, store, field.Template)
			if err != nil {
				return fmt.Errorf("failed to encode message property %s under %s: %w", fieldName, path, err)
			}
		}

	case tpl.Scalar != nil:
		var value interface{} // value to cast

		if ref == nil {
			value = tpl.Scalar.Default
		} else {
			value = ref.Value
		}

		casted := castType(tpl.Scalar.Type, value)
		if casted != "" {
			encoded.Add(path, casted)
		}

	case tpl.Enum != nil:
		if ref == nil {
			// no default value for nil. No reference => nothing to encode.
			break
		}

		value := tpl.Enum.Positions[*ref.Enum]
		casted := castType(types.Enum, value.Key)
		if casted != "" {
			encoded.Add(path, casted)
		}

	// repeated is described by a static template with a reference
	case tpl.Repeated != nil && ref != nil:
		item, err := tpl.Repeated.Template()

		if err != nil {
			return fmt.Errorf("failed to encode repeated property %s: %w", path, err)
		}

		// as item is the static template, it does not have its own Reference.
		// all repeated values are located in store by empty resource-path identifier.
		item.Reference = &specs.PropertyReference{Resource: "", Path: ""}

		for idx, store := range ref.Repeated {
			path := fmt.Sprintf("%s[%d]", path, idx)
			err = encode(encoded, path, store, item)

			if err != nil {
				return fmt.Errorf("failed to encode repeated property item %s: %w", path, err)
			}
		}

	// repeated does not have a static template but described "inline"
	case tpl.Repeated != nil && ref == nil:
		for idx, item := range tpl.Repeated {
			path := fmt.Sprintf("%s[%d]", path, idx)
			err := encode(encoded, path, store, item)

			if err != nil {
				return fmt.Errorf("failed to encode repeated property item %s: %w", path, err)
			}
		}
	}

	return nil
}

// Unmarshal the given www-form-urlencoded io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs references.Store) error {
	if manager.property == nil {
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
		if err := decodeElement(manager.resource, 0, strings.Split(key, "."), values, manager.property.Message, refs); err != nil {
			return err
		}
	}

	return nil
}

func decodeElement(resource string, pos int, path []string, values []string, schema specs.Message, refs references.Store) error {
	propName := path[pos]

	if schema == nil {
		return errNilSchema
	}

	prop := schema[propName]
	if prop == nil {
		return errUndefinedProperty(propName)
	}

	ref := &references.Reference{
		Path: prop.Path,
	}

	// TODO: implement nested repeated
	switch {
	case prop.Repeated != nil:
		for _, raw := range values {
			store := references.NewReferenceStore(0)

			switch {
			case prop.Template.Message != nil:
				if len(path) > pos+1 {
					if err := decodeElement(resource, pos+1, path, []string{raw}, prop.Message, store); err != nil {
						return err
					}
				}
			case prop.Template.Enum != nil:
				enum := prop.Enum.Keys[raw]
				if enum != nil {
					store.StoreEnum("", "", enum.Position)
				}
			default:
				value, err := types.DecodeFromString(raw, prop.Type())
				if err != nil {
					return err
				}

				store.StoreValue("", "", value)
			}

			ref.Append(store)
		}
	case prop.Message != nil:
		if len(path) > pos+1 {
			return decodeElement(resource, pos+1, path, values, prop.Message, refs)
		}
	case prop.Enum != nil:
		enum := prop.Enum.Keys[values[0]]
		if enum != nil {
			ref.Enum = &enum.Position
		}
	default:
		value, err := types.DecodeFromString(values[0], prop.Type())
		if err != nil {
			return err
		}

		ref.Value = value
	}

	refs.StoreReference(resource, ref)

	return nil
}
