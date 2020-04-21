package mock

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/jexia/maestro/pkg/definitions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"gopkg.in/yaml.v2"
)

// CollectionResolver returns the full mock collection
func CollectionResolver(path string) (*Collection, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	collection, err := UnmarshalFile(reader)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// SchemaResolver returns a new schema resolver for the given mock collection
func SchemaResolver(path string) definitions.SchemaResolver {
	return func(ctx instance.Context) ([]*specs.SchemaManifest, error) {
		reader, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		collection, err := UnmarshalFile(reader)
		if err != nil {
			return nil, err
		}

		return SchemaManifest(collection), nil
	}
}

// ServicesResolver returns a new service(s) resolver for the given mock collection
func ServicesResolver(path string) definitions.ServicesResolver {
	return func(ctx instance.Context) ([]*specs.ServicesManifest, error) {
		reader, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		collection, err := UnmarshalFile(reader)
		if err != nil {
			return nil, err
		}

		return ServiceManifest(collection), nil
	}
}

// UnmarshalFile attempts to parse the given Mock YAML file to intermediate resources.
func UnmarshalFile(reader io.Reader) (*Collection, error) {
	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	collection := Collection{}
	err = yaml.UnmarshalStrict(bb, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}
