package graphql

import (
	"testing"

	"github.com/graphql-go/graphql"
)

func TestSetFieldNil(t *testing.T) {
	_ = SetField("", nil, nil)
}

func TestSetFieldSimple(t *testing.T) {
	fields := graphql.Fields{}
	field := &graphql.Field{}

	err := SetField("message", fields, field)
	if err != nil {
		t.Fatal(err)
	}

	if fields["message"] != field {
		t.Fatalf("unexpected field %+v, expected %+v", fields["message"], field)
	}
}

func TestSetField(t *testing.T) {
	type test struct {
		path   string
		fields graphql.Fields
		field  *graphql.Field
	}

	tests := map[string]test{
		"unknown": {
			path:   "meta.time.unknown",
			fields: graphql.Fields{},
			field:  &graphql.Field{},
		},
		"simple": {
			path:   "message",
			fields: graphql.Fields{},
			field:  &graphql.Field{},
		},
		"nested": {
			path: "meta.value",
			fields: graphql.Fields{
				"meta": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name:   "meta",
						Fields: graphql.Fields{},
					}),
				},
			},
			field: &graphql.Field{},
		},
		"deep": {
			path: "meta.time.value",
			fields: graphql.Fields{
				"meta": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "meta",
						Fields: graphql.Fields{
							"time": &graphql.Field{
								Type: graphql.NewObject(graphql.ObjectConfig{
									Name:   "time",
									Fields: graphql.Fields{},
								}),
							},
						},
					}),
				},
			},
			field: &graphql.Field{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := SetField(test.path, test.fields, test.field)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSetFieldErr(t *testing.T) {
	type test struct {
		path   string
		fields graphql.Fields
		field  *graphql.Field
	}

	tests := map[string]test{
		"nested types": {
			path: "meta.value",
			fields: graphql.Fields{
				"meta": &graphql.Field{
					Type: graphql.Boolean,
				},
			},
			field: &graphql.Field{},
		},
		"deep types": {
			path: "meta.deep.value",
			fields: graphql.Fields{
				"meta": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "meta",
						Fields: graphql.Fields{
							"deep": &graphql.Field{
								Type: graphql.Boolean,
							},
						},
					}),
				},
			},
			field: &graphql.Field{},
		},
		"already defined": {
			path: "meta.value",
			fields: graphql.Fields{
				"meta": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "meta",
						Fields: graphql.Fields{
							"value": &graphql.Field{
								Type: graphql.Boolean,
							},
						},
					}),
				},
			},
			field: &graphql.Field{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := SetField(test.path, test.fields, test.field)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}
