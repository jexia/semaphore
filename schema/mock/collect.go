package mock

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/schema"
	"gopkg.in/yaml.v2"
)

// ErrResolver returns a new schema resolver which returns the given error when called
func ErrResolver(err error) schema.Resolver {
	return func(instance.Context, *schema.Store) error {
		return err
	}
}

// SchemaResolver returns a new schema resolver for the given mock collection
func SchemaResolver(path string) schema.Resolver {
	reader, err := os.Open(path)
	if err != nil {
		return ErrResolver(err)
	}

	collection, err := UnmarshalFile(reader)
	if err != nil {
		return ErrResolver(err)
	}

	return func(ctx instance.Context, schemas *schema.Store) error {
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
