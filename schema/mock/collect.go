package mock

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// UnmarshalFile attempts to parse the given Mock YAML file to intermediate resources.
func UnmarshalFile(path string) (*Collection, error) {
	bb, err := ioutil.ReadFile(path)
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
