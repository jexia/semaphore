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

func TestParseIntermediateStaticProperty(t *testing.T) {
	type test struct {
		template string
		expected *specs.Property
	}

	var tests = map[string]test{
		"array single": {
			template: `
				array = [
					"foo",
				]`,
			expected: &specs.Property{
				Name:  "array",
				Path:  "array",
				Label: labels.Optional,
				Template: specs.Template{
					Repeated: &specs.Repeated{
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "foo",
							},
						},
					},
				},
			},
		},
		"array reference": {
			template: `
				array = [
					"{{ input:message }}",
				]`,
			expected: &specs.Property{
				Name:  "array",
				Path:  "array",
				Label: labels.Optional,
				Template: specs.Template{
					Repeated: &specs.Repeated{
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: "input",
								Path:     "message",
							},
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
		"array of ints": {
			template: `
				array = [
					10,
					42
				]`,
			expected: &specs.Property{
				Name:  "array",
				Path:  "array",
				Label: labels.Optional,
				Template: specs.Template{
					Repeated: &specs.Repeated{
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.Int64, //
							},
						},
						Default: map[uint]*specs.Property{
							0: {
								Path:  "array", // ??? maybe use Template as well or even a custom struct with Reference and Value (interface{}) inside?
								Label: labels.Optional,
								Template: specs.Template{
									Scalar: &specs.Scalar{
										Type:    types.Int64,
										Default: int64(10),
									},
								},
							},
							1: {
								Path:  "array",
								Label: labels.Optional,
								Template: specs.Template{
									Scalar: &specs.Scalar{
										Type:    types.Int64,
										Default: int64(42),
									},
								},
							},
						},
					},
				},
			},
		},
		"array of objects": {
			template: `
				array = [
					{
						"action": "create",
					},
					{
						"action": "update",
					}
				]`,
			expected: &specs.Property{
				Name:  "array",
				Path:  "array",
				Label: labels.Optional,
				Template: specs.Template{
					Repeated: &specs.Repeated{
						Default: map[uint]*specs.Property{},
					},
				},

				Type: types.String,

				Nested: []*specs.Property{
					{
						Path:  "array",
						Type:  types.Message,
						Label: labels.Optional,
						Nested: []*specs.Property{
							{
								Name:    "action",
								Path:    "array.action",
								Type:    types.String,
								Label:   labels.Optional,
								Default: "create",
							},
						},
					},
					{
						Path:  "array",
						Type:  types.Message,
						Label: labels.Optional,
						Nested: []*specs.Property{
							{
								Name:    "action",
								Path:    "array.action",
								Type:    types.String,
								Label:   labels.Optional,
								Default: "update",
							},
						},
					},
				},
			},
		},
		"array complex": {
			template: `
				array = [
					"foo",
					42,
					{
						"id": "{{ input:id }}",
					}
				]`,
			expected: &specs.Property{
				Name:  "array",
				Path:  "array",
				Type:  types.String,
				Label: labels.Repeated,
				Nested: []*specs.Property{
					{
						Path:    "array",
						Type:    types.String,
						Label:   labels.Optional,
						Default: "foo",
					},
					{
						Path:    "array",
						Type:    types.Int64,
						Label:   labels.Optional,
						Default: int64(42),
					},
					{
						Path:  "array",
						Type:  types.Message,
						Label: labels.Optional,
						Nested: []*specs.Property{
							{
								Name: "id",
								Path: "array.id",
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "id",
								},
							},
						},
					},
				},
			},
		},
		"object complex": {
			template: `
				object = {
					"message": "hello world",
					"meta": {
						"id": "{{ input:id }}",
						"foo": 42
					}
				}`,
			expected: &specs.Property{
				Name:  "object",
				Path:  "object",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: []*specs.Property{
					{
						Name:    "message",
						Path:    "object.message",
						Type:    types.String,
						Label:   labels.Optional,
						Default: "hello world",
					},
					{
						Name:  "meta",
						Path:  "object.meta",
						Type:  types.Message,
						Label: labels.Optional,
						Nested: []*specs.Property{
							{
								Name: "id",
								Path: "object.meta.id",
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "id",
								},
							},
							{
								Name:    "foo",
								Path:    "object.meta.foo",
								Label:   labels.Optional,
								Type:    types.Int64,
								Default: int64(42),
							},
						},
					},
				},
			},
		},
		"object array": {
			template: `
				object = {
					"message": [
						"hello world",
						"{{ input:message }}"
					],
				}`,
			expected: &specs.Property{
				Name:  "object",
				Path:  "object",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: []*specs.Property{
					{
						Name:  "message",
						Path:  "object.message",
						Type:  types.String,
						Label: labels.Repeated,
						Nested: []*specs.Property{
							{
								Path:    "object.message",
								Type:    types.String,
								Label:   labels.Optional,
								Default: "hello world",
							},
							{
								Path: "object.message",
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "message",
								},
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
				attr   = parseTestAttribute(t, test.template)
				ctx    = logger.WithLogger(broker.NewBackground())
				val, _ = attr.Expr.Value(nil)
			)

			result, err := ParseIntermediateProperty(ctx, "", attr, val)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			ValidateProperties(t, result, test.expected)
		})
	}
}

func ValidateProperties(t *testing.T, result *specs.Property, expected *specs.Property) {
	if result.Path != expected.Path {
		t.Errorf("property path %q was expected to be %q", result.Path, expected.Path)
	}

	if result.Name != expected.Name {
		t.Errorf("property (%q) name %q was expected to be %q", result.Path, result.Name, expected.Name)
	}

	if result.Label != expected.Label {
		t.Errorf("property (%q) label %q was expected to be %q", result.Path, result.Label, expected.Label)
	}

	if result.Type != expected.Type {
		t.Errorf("property (%q) type %q was expected to be %q", result.Path, result.Type, expected.Type)
	}

	if !reflect.DeepEqual(result.Reference, expected.Reference) {
		t.Errorf("property (%q) reference \"%v\" was expected to be \"%v\"", result.Path, result.Reference, expected.Reference)
	}

	if !reflect.DeepEqual(result.Default, expected.Default) {
		t.Errorf("property (%q) default value \"%v\" was expected to be \"%v\"", result.Path, result.Default, expected.Default)
	}

	if len(result.Nested) != len(expected.Nested) {
		t.Fatalf("unexpected repeated %+v, expected %+v", result.Nested, expected.Nested)
	}

	for index, schema := range expected.Nested {
		// if array order matters
		nested := result.Nested[index]
		if expected.Type == types.Message {
			nested = result.Nested.Get(schema.Name)
		}

		if nested == nil {
			t.Fatalf("was expected to contain nested property %q inside %q", schema.Name, result.Path)
		}

		ValidateProperties(t, nested, schema)
	}
}
