package specs

import (
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestCompareRepeated(t *testing.T) {
	repeatedA := Repeated{
		{
			Scalar: &Scalar{
				Type: types.Int32,
			},
		},
	}

	repeatedB := Repeated{
		{
			Scalar: &Scalar{
				Type: types.String,
			},
		},
	}

	type args struct {
		expected Repeated
		given    Repeated
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"should not match",
			args{repeatedA, repeatedB},
			true,
		},

		{
			"should match",
			args{repeatedA, repeatedA},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.given.Compare(tt.args.expected); (err != nil) != tt.wantErr {
				t.Errorf("CompareRepeated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepeatedTemplate(t *testing.T) {
	type test struct {
		repeated Repeated
		expected Template
		isError  bool
	}

	var (
		foo = Template{
			Scalar: &Scalar{
				Type:    types.String,
				Default: "foo",
			},
		}

		bar = Template{
			Scalar: &Scalar{
				Type:    types.String,
				Default: "bar",
			},
		}

		baz = Template{
			Scalar: &Scalar{
				Type:    types.String,
				Default: "baz",
			},
		}

		fooBarMsg = Message{
			"foo": &Property{Template: foo},
			"bar": &Property{Template: bar},
		}

		objFooBar = Template{
			Message: fooBarMsg,
		}

		barBazMsg = Message{
			"bar": &Property{Template: bar},
			"baz": &Property{Template: baz},
		}

		objBarBaz = Template{
			Message: barBazMsg,
		}

		tests = map[string]test{
			"strings": {
				repeated: Repeated{foo, bar, baz},
				expected: Template{
					Scalar: &Scalar{
						Type: types.String,
					},
				},
			},
			"objects": {
				repeated: Repeated{objFooBar, objFooBar, objFooBar},
				expected: Template{
					Message: fooBarMsg,
				},
			},
			"type mismatch": {
				repeated: Repeated{objFooBar, objBarBaz},
				isError:  true,
			},
		}
	)

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			actual, err := test.repeated.Template()

			if test.isError && err == nil {
				t.Error("error was expected")
			}

			if !test.isError && err != nil {
				t.Errorf("unexpected error '%s'", err)
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("template %+v was expected to be %+v", actual, test.expected)
			}
		})
	}
}
