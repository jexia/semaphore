package mock

import (
	"io"
	"io/ioutil"

	"github.com/jexia/maestro/schema"
	"gopkg.in/yaml.v2"
)

// Read reads the given io reader and returns a schema resolver
func Read(reader io.Reader) (schema.Resolver, error) {
	collection, err := UnmarshalFile(reader)
	if err != nil {
		return nil, err
	}

	return SchemaResolver(collection), nil
}

// SchemaResolver returns a new schema resolver for the given mock collection
func SchemaResolver(collection schema.Collection) schema.Resolver {
	return func(schemas *schema.Store) error {
		schemas.Add(collection)
		return nil
	}
}

// UnmarshalFile attempts to parse the given Mock YAML file to intermediate resources.
func UnmarshalFile(reader io.Reader) (*Collection, error) {
	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	collection := Collection{}
	err = yaml.Unmarshal(bb, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}
