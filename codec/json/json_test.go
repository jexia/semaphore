package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs"
)

func NewMock(t *testing.T) specs.Object {
	path, err := filepath.Abs("./tests/logger.yaml")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := os.Open(path)
	collection, err := mock.UnmarshalFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := maestro.New(maestro.WithPath("./tests", false), maestro.WithSchemaCollection(collection))
	if err != nil {
		t.Fatal(err)
	}

	return manifest.Flows[0].GetCalls()[0].Request
}

func TestMarshal(t *testing.T) {
	specs := NewMock(t)
	manager, err := New("input", nil, specs)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]map[string]interface{}{
		"simple": map[string]interface{}{
			"message": "some message",
			"nested":  map[string]interface{}{},
		},
		"nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "some message",
			},
		},
		"complex": map[string]interface{}{
			"message": "hello world",
			"nested": map[string]interface{}{
				"value": "nested value",
			},
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
			},
		},
	}

	for key, input := range tests {
		t.Run(key, func(t *testing.T) {
			inputAsJSON, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			refs := refs.NewStore(len(input))
			refs.StoreValues("input", "", input)

			reader, err := manager.Marshal(refs)
			if err != nil {
				t.Fatal(err)
			}

			responseAsJSON, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			result := map[string]interface{}{}
			err = json.Unmarshal(responseAsJSON, &result)
			if err != nil {
				t.Fatal(err)
			}

			expected := map[string]interface{}{}
			err = json.Unmarshal(inputAsJSON, &expected)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expected, result) {
				t.Errorf("unexpected response %s, expected %s", string(responseAsJSON), string(inputAsJSON))
			}
		})
	}
}
