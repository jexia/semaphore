package graphql

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

func TestNewObjectNil(t *testing.T) {
	NewObject("", "", specs.Template{})
}

func TestNewObject(t *testing.T) {
	type test struct {
		template specs.Template
	}

	tests := map[string]test{
		"scalar": {
			template: specs.Template{
				Message: specs.Message{
					"scalar": &specs.Property{
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
		"message": {
			template: specs.Template{
				Message: specs.Message{
					"nested": &specs.Property{
						Template: specs.Template{
							Message: specs.Message{
								"value": &specs.Property{
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"repeated": {
			template: specs.Template{
				Message: specs.Message{
					"repeated": &specs.Property{
						Template: specs.Template{
							Repeated: specs.Repeated{
								specs.Template{
									Scalar: &specs.Scalar{
										Type: types.String,
									},
								},
							},
						},
					},
				},
			},
		},
		"enum": {
			template: specs.Template{
				Message: specs.Message{
					"enum": &specs.Property{
						Template: specs.Template{
							Enum: &specs.Enum{
								Keys: map[string]*specs.EnumValue{
									"UNKOWN": {
										Key:      "UNKOWN",
										Position: 0,
									},
								},
								Positions: map[int32]*specs.EnumValue{
									0: {
										Key:      "UNKOWN",
										Position: 0,
									},
								},
							},
							Reference: &specs.PropertyReference{
								Resource: "input",
								Path:     "repeated.status",
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewObject("", "", test.template)
			if err != nil {
				t.Fatal(err)
			}

			// TODO: validate object
		})
	}
}

func TestNewSchemaObject(t *testing.T) {
	type test struct {
		template specs.Template
	}

	tests := map[string]test{
		"scalar": {
			template: specs.Template{
				Message: specs.Message{
					"scalar": &specs.Property{
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
		"message": {
			template: specs.Template{
				Message: specs.Message{
					"nested": &specs.Property{
						Template: specs.Template{
							Message: specs.Message{
								"value": &specs.Property{
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"repeated": {
			template: specs.Template{
				Message: specs.Message{
					"repeated": &specs.Property{
						Template: specs.Template{
							Repeated: specs.Repeated{
								specs.Template{
									Scalar: &specs.Scalar{
										Type: types.String,
									},
								},
							},
						},
					},
				},
			},
		},
		"enum": {
			template: specs.Template{
				Message: specs.Message{
					"enum": &specs.Property{
						Template: specs.Template{
							Enum: &specs.Enum{
								Keys: map[string]*specs.EnumValue{
									"UNKOWN": {
										Key:      "UNKOWN",
										Position: 0,
									},
								},
								Positions: map[int32]*specs.EnumValue{
									0: {
										Key:      "UNKOWN",
										Position: 0,
									},
								},
							},
							Reference: &specs.PropertyReference{
								Resource: "input",
								Path:     "repeated.status",
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			object := &transport.Object{
				Definition: &specs.ParameterMap{
					Property: &specs.Property{
						Template: test.template,
					},
				},
			}

			objects := NewObjects()
			_, err := NewSchemaObject(objects, "", object)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
