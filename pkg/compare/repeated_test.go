package compare

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestCompareRepeated(t *testing.T) {
	repeatedA := &specs.Repeated{
		Property: &specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{Type: types.Int32},
			},
		},
	}

	repeatedB := &specs.Repeated{
		Property: &specs.Property{
			Template: specs.Template{
				Scalar: &specs.Scalar{Type: types.String},
			},
		},
	}

	type args struct {
		expected *specs.Repeated
		given    *specs.Repeated
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
			if err := CompareRepeated(tt.args.expected, tt.args.given); (err != nil) != tt.wantErr {
				t.Errorf("CompareRepeated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
