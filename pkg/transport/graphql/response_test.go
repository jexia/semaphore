package graphql

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestResponseObjectNil(t *testing.T) {
	ResponseObject(nil, nil, nil)
}

func TestResponseObjectInvalidType(t *testing.T) {
	store := references.NewStore(0)
	tracker := references.NewTracker()
	property := &specs.Property{
		Template: specs.Template{
			Scalar: &specs.Scalar{},
		},
	}

	_, err := ResponseObject(property, store, tracker)
	if err != ErrInvalidObject {
		t.Fatalf("unexpected err %s, expected %s", err, ErrInvalidObject)
	}
}

func TestResponseObject(t *testing.T) {
	type populate func(t *testing.T, store references.Store)
	type test struct {
		property *specs.Property
		populate populate
		expected map[string]interface{}
	}

	tests := map[string]test{
		"empty": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{},
				},
			},
			populate: func(t *testing.T, store references.Store) {},
			expected: map[string]interface{}{},
		},
		"object": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"name": &specs.Property{
							Name: "name",
							Path: "name",
							Template: specs.Template{
								Scalar: &specs.Scalar{
									Type: types.String,
								},
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "name",
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {
				store.Store(template.ResourcePath(template.InputResource, "name"), &references.Reference{Value: "John Doe"})
			},
			expected: map[string]interface{}{
				"name": "John Doe",
			},
		},
		"nested object": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"nested": &specs.Property{
							Name: "nested",
							Path: "nested",
							Template: specs.Template{
								Message: specs.Message{
									"name": &specs.Property{
										Name: "name",
										Path: "nested.name",
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Type: types.String,
											},
											Reference: &specs.PropertyReference{
												Resource: "input",
												Path:     "name",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {
				store.Store(template.ResourcePath(template.InputResource, "name"), &references.Reference{Value: "John Doe"})
			},
			expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"name": "John Doe",
				},
			},
		},
		"repeated object": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"repeated": &specs.Property{
							Name: "repeated",
							Path: "repeated",
							Template: specs.Template{
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "repeated",
								},
								Repeated: specs.Repeated{
									specs.Template{
										Message: specs.Message{
											"name": &specs.Property{
												Name: "name",
												Path: "repeated.name",
												Template: specs.Template{
													Scalar: &specs.Scalar{
														Type: types.String,
													},
													Reference: &specs.PropertyReference{
														Resource: "input",
														Path:     "repeated.name",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {
				tracker := references.NewTracker()
				references.StoreValues(store, tracker, template.ResourcePath(template.InputResource), map[string]interface{}{
					"repeated": []map[string]interface{}{
						{
							"name": "John Doe",
						},
					},
				})

			},
			expected: map[string]interface{}{
				"repeated": []interface{}{
					map[string]interface{}{
						"name": "John Doe",
					},
				},
			},
		},
		"object nil keys": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"name": &specs.Property{
							Name: "name",
							Path: "name",
							Template: specs.Template{
								Scalar: &specs.Scalar{
									Type: types.String,
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {},
			expected: map[string]interface{}{},
		},
		"object nil reference": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"name": &specs.Property{
							Name: "name",
							Path: "name",
							Template: specs.Template{
								Scalar: &specs.Scalar{
									Type: types.String,
								},
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "name",
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {},
			expected: map[string]interface{}{},
		},
		"enum": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"status": &specs.Property{
							Name: "status",
							Path: "status",
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
									Path:     "status",
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {
				enum := int32(0)
				store.Store(template.ResourcePath(template.InputResource, "status"), &references.Reference{Enum: &enum})
			},
			expected: map[string]interface{}{
				"status": "UNKOWN",
			},
		},
		"enum nil reference": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"status": &specs.Property{
							Name: "status",
							Path: "status",
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
									Path:     "status",
								},
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {},
			expected: map[string]interface{}{},
		},
		"enum nil value": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"status": &specs.Property{
							Name: "status",
							Path: "status",
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
							},
						},
					},
				},
			},
			populate: func(t *testing.T, store references.Store) {},
			expected: map[string]interface{}{},
		},
		"repeated enum": {
			property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"repeated": &specs.Property{
							Name: "repeated",
							Path: "repeated",
							Template: specs.Template{
								Reference: &specs.PropertyReference{
									Resource: "input",
									Path:     "repeated",
								},
								Repeated: specs.Repeated{
									specs.Template{
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
				},
			},
			populate: func(t *testing.T, store references.Store) {
				tracker := references.NewTracker()
				references.StoreValues(store, tracker, template.ResourcePath(template.InputResource), map[string]interface{}{
					"repeated": []map[string]interface{}{
						{
							"status": references.Enum("UNKOWN", 0),
						},
					},
				})
			},
			expected: map[string]interface{}{
				"repeated": []interface{}{
					"UNKOWN",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			store := references.NewStore(0)
			tracker := references.NewTracker()

			test.populate(t, store)

			response, err := ResponseObject(test.property, store, tracker)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(response, test.expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
