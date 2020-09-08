package hcl

import (
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

const (
	hclV2array = `array = [
		"foo",
		42,
		{
			"id": "{{ getter:output }}",
		}
	]`

	hclV2object = `object = {
		"message": "hello world",
		"meta": {
			"id": "{{ getter:output }}",
			"foo": 42
		}
	}`
)

func parseTestAttribute(t *testing.T, config string) *hcl.Attribute {
	bb := []byte(os.ExpandEnv(config))

	file, diags := hclsyntax.ParseConfig(bb, "test", hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("unexpected error: %s", diags.Error())
	}

	var parameterMap ParameterMap

	diags = gohcl.DecodeBody(file.Body, nil, &parameterMap)
	if diags.HasErrors() {
		t.Fatalf("unexpected error: %s", diags.Error())
	}

	attrs, _ := parameterMap.Properties.JustAttributes()
	for _, attr := range attrs {
		return attr
	}

	t.Fatal("must contain at least one attribute")

	return nil
}

type output struct {
	name         string
	path         string
	dataType     types.Type
	label        labels.Label
	nested       map[string]output
	repeated     []output
	reference    *specs.PropertyReference
	defaultValue interface{}
}

func TestParseIntermediateProperty(t *testing.T) {
	type test struct {
		input  string
		output *output
	}

	var tests = map[string]test{
		"HCL2 array": {
			input: hclV2array,
			output: &output{
				name:     "array",
				path:     "array",
				dataType: types.String, // TODO: get rid of hardcoded type
				label:    labels.Repeated,
				repeated: []output{
					{
						path:         "array",
						dataType:     types.String,
						label:        labels.Optional,
						defaultValue: "foo",
					},
					{
						path:         "array",
						dataType:     types.Int64,
						label:        labels.Optional,
						defaultValue: int64(42),
					},
					{
						path:     "array",
						dataType: types.Message,
						label:    labels.Optional,
						nested: map[string]output{
							"id": {
								name: "id",
								path: "array.id",
								reference: &specs.PropertyReference{
									Resource: "getter",
									Path:     "output",
								},
							},
						},
					},
				},
			},
		},
		"HCL2 object": {
			input: hclV2object,
			output: &output{
				name:     "object",
				path:     "object",
				dataType: types.Message,
				label:    labels.Optional,
				nested: map[string]output{
					"message": {
						name:         "message",
						path:         "object.message",
						dataType:     types.String,
						label:        labels.Optional,
						defaultValue: "hello world",
					},
					"meta": {
						name:     "meta",
						path:     "object.meta",
						dataType: types.Message,
						label:    labels.Optional,
						nested: map[string]output{
							"id": {
								name: "id",
								path: "object.meta.id",
								reference: &specs.PropertyReference{
									Resource: "getter",
									Path:     "output",
								},
							},
							"foo": {
								name:         "foo",
								path:         "object.meta.foo",
								label:        labels.Optional,
								dataType:     types.Int64,
								defaultValue: int64(42),
							},
						},
					},
				},
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			var (
				attr   = parseTestAttribute(t, test.input)
				ctx    = logger.WithLogger(broker.NewBackground())
				val, _ = attr.Expr.Value(nil)
			)

			prop, err := ParseIntermediateProperty(ctx, "", attr, val)
			if err != nil {
				t.Errorf("unexpected error: %s", err)

				return
			}

			assertProperties(t, prop, test.output)
		})
	}
}

func assertProperties(t *testing.T, prop *specs.Property, output *output) {
	if output == nil {
		return
	}

	if prop.Name != output.name {
		t.Errorf("property name %q was expected to be %q", prop.Name, output.name)
	}

	if prop.Path != output.path {
		t.Errorf("property path %q was expected to be %q", prop.Path, output.path)
	}

	if prop.Label != output.label {
		t.Errorf("property label %q was expected to be %q", prop.Label, output.label)
	}

	if prop.Type != output.dataType {
		t.Errorf("property type %q was expected to be %q", prop.Type, output.dataType)
	}

	if !reflect.DeepEqual(prop.Reference, output.reference) {
		t.Errorf("property reference \"%v\" was expected to be \"%v\"", prop.Reference, output.reference)
	}

	if !reflect.DeepEqual(prop.Default, output.defaultValue) {
		t.Errorf("default value \"%v\" was expected to be \"%v\"", prop.Default, output.defaultValue)
	}

	if actual, expected := len(prop.Repeated), len(output.repeated); actual != expected {
		t.Errorf("got %d repeated values, expected %d", actual, expected)

		return
	}

	for index, repeated := range prop.Repeated {
		assertProperties(t, repeated, &output.repeated[index])
	}

	if actual, expected := len(prop.Nested), len(output.nested); actual != expected {
		t.Errorf("got %d nested values, expected %d", actual, expected)

		return
	}

	for key, expected := range output.nested {
		actual, ok := prop.Nested[key]
		if !ok {
			t.Errorf("was expected to contain nested property %q", key)
		}

		assertProperties(t, actual, &expected)
	}
}
