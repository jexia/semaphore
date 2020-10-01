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
					Repeated: specs.Repeated{
						specs.Template{
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
					Repeated: specs.Repeated{
						specs.Template{
							Reference: &specs.PropertyReference{
								Resource: "input",
								Path:     "message",
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
					Repeated: specs.Repeated{
						specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.Int64,
								Default: int64(10),
							},
						},
						specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.Int64,
								Default: int64(42),
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
					Repeated: specs.Repeated{
						specs.Template{
							Message: specs.Message{
								"action": {
									Name:  "action",
									Path:  "array.action",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type:    types.String,
											Default: "create",
										},
									},
								},
							},
						},
						specs.Template{
							Message: specs.Message{
								"action": {
									Name:  "action",
									Path:  "array.action",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type:    types.String,
											Default: "update",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"array complex": {
			template: `
				array = [
					{
						"id": "{{ input:id }}",
					},
					"{{ input:name }}"
				]`,
			expected: &specs.Property{
				Name:  "array",
				Path:  "array",
				Label: labels.Optional,
				Template: specs.Template{
					Repeated: specs.Repeated{
						specs.Template{
							Message: specs.Message{
								"id": {
									Name: "id",
									Path: "array.id",
									Template: specs.Template{
										Reference: &specs.PropertyReference{
											Resource: "input",
											Path:     "id",
										},
									},
								},
							},
						},
						specs.Template{
							Reference: &specs.PropertyReference{
								Resource: "input",
								Path:     "name",
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
				Label: labels.Optional,
				Template: specs.Template{
					Message: specs.Message{
						"message": {
							Name:  "message",
							Path:  "object.message",
							Label: labels.Optional,
							Template: specs.Template{
								Repeated: specs.Repeated{
									specs.Template{
										Scalar: &specs.Scalar{
											Type:    types.String,
											Default: "hello world",
										},
									},
									specs.Template{
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
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			attr := parseTestAttribute(t, test.template)
			ctx := logger.WithLogger(broker.NewBackground())
			val, _ := attr.Expr.Value(nil)

			result, err := ParseIntermediateProperty(ctx, "", attr, val)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			ValidateProperty(t, result, test.expected)
		})
	}
}

func ValidateProperty(t *testing.T, result *specs.Property, expected *specs.Property) {
	if result.Path != expected.Path {
		t.Errorf("property path %q was expected to be %q", result.Path, expected.Path)
	}

	if result.Name != expected.Name {
		t.Errorf("property (%q) name %q was expected to be %q", result.Path, result.Name, expected.Name)
	}

	if result.Label != expected.Label {
		t.Errorf("property (%q) label %q was expected to be %q", result.Path, result.Label, expected.Label)
	}

	ValidateTemplate(t, result.Path, result.Template, expected.Template)
}

func ValidateTemplate(t *testing.T, path string, result, expected specs.Template) {
	if expected.Scalar != nil && result.Scalar == nil {
		t.Errorf("property (%q) scalar was expected to be set", path)
	}

	if result.Scalar != nil && result.Scalar.Type != expected.Scalar.Type {
		t.Errorf("property (%q) type %q was expected to be %q", path, result.Scalar.Type, expected.Scalar.Type)
	}

	if result.Scalar != nil && !reflect.DeepEqual(result.Scalar.Default, expected.Scalar.Default) {
		t.Errorf("property (%q) default value \"%v\" was expected to be \"%v\"", path, result.Scalar.Default, expected.Scalar.Default)
	}

	if !reflect.DeepEqual(result.Reference, expected.Reference) {
		t.Errorf("property (%q) reference \"%v\" was expected to be \"%v\"", path, result.Reference, expected.Reference)
	}

	if len(result.Message) != len(expected.Message) {
		t.Fatalf("unexpected message %+v, expected %+v", result.Message, expected.Message)
	}

	if len(result.Repeated) != len(expected.Repeated) {
		t.Fatalf("unexpected repeated %+v, expected %+v", result.Repeated, expected.Repeated)
	}

	for key, schema := range expected.Message {
		nested := result.Message[key]
		if nested == nil {
			t.Fatalf("was expected to contain message property %q inside %q", schema.Name, path)
		}

		ValidateProperty(t, nested, schema)
	}

	for index, schema := range expected.Repeated {
		nested := result.Repeated[index]
		ValidateTemplate(t, path, nested, schema)
	}
}
