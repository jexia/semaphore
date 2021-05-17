package references

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/specs"
)

func TestReferencedCollectionSet(t *testing.T) {
	tests := []struct {
		name string
		want map[string]struct{}
	}{
		{
			name: "return set",
			want: map[string]struct{}{"some": {}, "some.random": {}, "some.random.path": {}},
		},
	}

	for _, test := range tests {
		r := ReferencedCollection{}
		r.Set("some.random.path")
		for want := range test.want {
			if _, ok := r[want]; !ok {
				t.Errorf("not set %s", want)
			}
		}
	}
}

func TestReferencedCollectionHas(t *testing.T) {
	tests := []struct {
		name string
		want map[string]struct{}
	}{
		{
			name: "return set",
			want: map[string]struct{}{"some": {}, "some.random": {}, "some.random.path": {}},
		},
	}

	for _, test := range tests {
		r := ReferencedCollection{}
		r.Set("some.random.path")
		for want := range test.want {
			if !r.Has(want) {
				t.Errorf("not set %s", want)
			}
		}
	}
}

func TestReferencedParameterMapPaths(t *testing.T) {
	tests := []struct {
		name   string
		want   map[string]struct{}
		params *specs.ParameterMap
		ref    ReferencedCollection
	}{
		{
			name: "repeated not nil",
			want: map[string]struct{}{
				"some.property": {},
			},
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Template: specs.Template{
						Repeated: []specs.Template{
							{
								Message: specs.Message{
									"some.property": &specs.Property{
										Path: "some.property",
									},
								},
							},
							{
								Message: specs.Message{
									"some.other.property": &specs.Property{
										Path: "some.other.property",
									},
								},
							},
						},
					},
				},
			},
			ref: ReferencedCollection{"some.property": {}, "": {}},
		},
		{
			name: "repeated nil",
			want: map[string]struct{}{
				"some.property":       {},
				"some.other.property": {},
			},
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Template: specs.Template{
						Repeated: []specs.Template{
							{
								Message: specs.Message{
									"some.property": &specs.Property{
										Path: "some.property",
									},
								},
							},
							{
								Message: specs.Message{
									"some.other.property": &specs.Property{
										Path: "some.other.property",
									},
								},
							},
						},
					},
				},
			},
			ref: ReferencedCollection{"some.property": {}},
		},
	}
	for _, test := range tests {
		newSpecMap := ReferencedParameterMapPaths(test.ref, test.params)
		for _, template := range newSpecMap.Property.Template.Repeated {
			for key := range template.Message {
				if _, ok := test.want[key]; !ok {
					t.Errorf("expected key not found %s", key)
				}
			}
		}
	}
}

func TestReferencedResourcePaths(t *testing.T) {
	tests := []struct {
		name     string
		want     map[string]struct{}
		flow     specs.Flow
		resource string
	}{
		{
			name:     "with data",
			resource: "user",
			want: map[string]struct{}{
				"some":               {},
				"some.property":      {},
				"some.property.user": {},
			},
			flow: specs.Flow{
				Name: "get_data",
				Nodes: specs.NodeList{
					{
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Header: map[string]*specs.Property{"Content-Type": nil},
							},
							Request: &specs.ParameterMap{
								Params: map[string]*specs.Property{"some.param": nil},
								Stack:  map[string]*specs.Property{"stack": nil},
								Header: map[string]*specs.Property{"Authorization": {
									Template: specs.Template{
										Repeated: []specs.Template{
											{
												Message: specs.Message{
													"some.property": &specs.Property{
														Template: specs.Template{
															Reference: &specs.PropertyReference{
																Resource: "user",
																Path:     "some.property.user",
															},
														},
													},
												},
											},
										},
									},
								}},
							},
						},
						Condition: &specs.Condition{},
						Intermediate: &specs.ParameterMap{
							Params:   map[string]*specs.Property{"username": {}},
							Property: &specs.Property{},
							Stack:    map[string]*specs.Property{"username": {}},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		collection := ReferencedResourcePaths(&test.flow, "user")
		for key := range test.want {
			if _, ok := collection[key]; !ok {
				t.Errorf("expected key not found %s", key)
			}
		}
	}
}
