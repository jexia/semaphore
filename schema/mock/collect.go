package mock

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
