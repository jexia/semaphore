package graphql

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestNewArgsNil(t *testing.T) {
	NewArgs(nil)
}

func TestNewArgsNilProperty(t *testing.T) {
	NewArgs(&specs.ParameterMap{})
}

func TestNewArgsUnexpectedType(t *testing.T) {
	props := &specs.ParameterMap{
		Property: &specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		},
	}

	_, err := NewArgs(props)
	_, is := err.(ErrUnexpectedType)
	if !is {
		t.Fatalf("unexpected error %+v, exoected ErrUnexpectedType", err)
	}
}

func TestNewArgsEmptyMessage(t *testing.T) {
	props := &specs.ParameterMap{
		Property: &specs.Property{
			Template: specs.Template{
				Message: specs.Message{},
			},
		},
	}

	field, err := NewArgs(props)
	if err != nil {
		t.Fatal(err)
	}

	if len(field) != 0 {
		t.Fatalf("unexpected field %+v, expected empty field", field)
	}
}

func TestNewArgs(t *testing.T) {
	tests := map[string]*specs.Property{
		"message": {
			Template: specs.Template{
				Message: specs.Message{},
			},
		},
		"repeated": {
			Template: specs.Template{
				Repeated: specs.Repeated{},
			},
		},
		"enum": {
			Template: specs.Template{
				Enum: &specs.Enum{
					Keys: map[string]*specs.EnumValue{
						"mock": {},
					},
				},
			},
		},
		"scalar": {
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			// ensure input name
			input.Name = name

			props := &specs.ParameterMap{
				Property: &specs.Property{
					Name: "mock",
					Template: specs.Template{
						Message: specs.Message{
							name: input,
						},
					},
				},
			}

			field, err := NewArgs(props)
			if err != nil {
				t.Fatal(err)
			}

			if len(field) != 1 {
				t.Fatalf("unexpected field %+v, expected single field", field)
			}
		})
	}
}

func TestNewInputArgsObjectNil(t *testing.T) {
	NewInputArgObject(nil)
}

func TestNewInputArgsObjectNilProperty(t *testing.T) {
	NewInputArgObject(&specs.Property{})
}

func TestNewInputArgsObjectUnexpectedType(t *testing.T) {
	prop := &specs.Property{}

	_, err := NewInputArgObject(prop)
	_, is := err.(ErrUnexpectedType)
	if !is {
		t.Fatalf("unexpected error %+v, exoected ErrUnexpectedType", err)
	}
}

func TestNewInputArgsObjectNoName(t *testing.T) {
	prop := &specs.Property{
		Template: specs.Template{
			Message: specs.Message{},
		},
	}

	expected := "Type must be named."
	object, err := NewInputArgObject(prop)
	if err.Error() != expected {
		t.Fatalf("unexpected err %+v, expected %s", err, expected)
	}

	if object != nil {
		t.Fatalf("unexpected object %+v, expected nil object", object)
	}
}

func TestNewInputArgsObjectEmptyMessage(t *testing.T) {
	prop := &specs.Property{
		Name: "mock",
		Template: specs.Template{
			Message: specs.Message{},
		},
	}

	object, err := NewInputArgObject(prop)
	if err != nil {
		t.Fatal(err)
	}

	if len(object.Fields()) != 0 {
		t.Fatalf("unexpected field %+v, expected empty field", object.Fields())
	}
}

func TestNewInputArgsObject(t *testing.T) {
	tests := map[string]*specs.Property{
		"message": {
			Template: specs.Template{
				Message: specs.Message{
					"value": &specs.Property{
						Name: "value",
						Path: "value",
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
		"repeated": {
			Template: specs.Template{
				Repeated: specs.Repeated{},
			},
		},
		"enum": {
			Template: specs.Template{
				Enum: &specs.Enum{
					Keys: map[string]*specs.EnumValue{
						"mock": {},
					},
				},
			},
		},
		"scalar": {
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			// ensure input name
			input.Name = name

			prop := &specs.Property{
				Name: "mock",
				Path: "mock",
				Template: specs.Template{
					Message: specs.Message{
						name: input,
					},
				},
			}

			object, err := NewInputArgObject(prop)
			if err != nil {
				t.Fatal(err)
			}

			if len(object.Fields()) != 1 {
				t.Fatalf("unexpected field %+v, expected single field", object.Fields())
			}
		})
	}
}
