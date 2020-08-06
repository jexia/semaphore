package mock

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
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
func SchemaResolver(path string) providers.SchemaResolver {
	return func(ctx instance.Context) (specs.Schemas, error) {
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
func ServicesResolver(path string) providers.ServicesResolver {
	return func(ctx instance.Context) (specs.ServiceList, error) {
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
