package formatter

import "encoding/json"

type JSON struct{}

func (JSON) String() string { return "json" }

func (JSON) Format(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
