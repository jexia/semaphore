package mock

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
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
	return func(ctx *broker.Context) (specs.Schemas, error) {
		reader, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		collection, err := UnmarshalFile(reader)
		if err != nil {
			return nil, err
		}

		return collection.Properties, nil
	}
}

// ServicesResolver returns a new service(s) resolver for the given mock collection
func ServicesResolver(path string) providers.ServicesResolver {
	return func(ctx *broker.Context) (specs.ServiceList, error) {
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

	DefinePropertyPaths(&collection)
	return &collection, nil
}

// DefinePropertyPaths defines all property paths inside the given collection
func DefinePropertyPaths(collection *Collection) {
	for _, property := range collection.Properties {
		definePath("", property)
	}
}

func definePath(path string, property *specs.Property) {
	fqpath := template.JoinPath(path, property.Name)
	property.Path = fqpath
	walkTemplate(fqpath, &property.Template)
}

func walkTemplate(path string, template *specs.Template) {
	if template.Message != nil {
		for key, property := range template.Message {
			property.Name = key
			definePath(path, property)
		}
	}

	if template.Repeated != nil {
		for _, repeated := range template.Repeated {
			walkTemplate(path, &repeated)
		}
	}
}
