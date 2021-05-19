package specs

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

func TestTemplate_Type(t *testing.T) {
	type fields struct {
		Scalar   *Scalar
		Enum     *Enum
		Repeated Repeated
		Message  Message
		OneOf    OneOf
	}
	tests := []struct {
		name   string
		fields fields
		want   types.Type
	}{
		{
			"return scalar type",
			fields{Scalar: &Scalar{Type: types.Int32}},
			types.Int32,
		},
		{
			"return enum",
			fields{Enum: &Enum{}},
			types.Enum,
		},
		{
			"return array",
			fields{Repeated: Repeated{}},
			types.Array,
		},
		{
			"return message",
			fields{Message: Message{}},
			types.Message,
		},
		{
			"return oneOf",
			fields{OneOf: OneOf{}},
			types.OneOf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := Template{
				Scalar:   tt.fields.Scalar,
				Enum:     tt.fields.Enum,
				Repeated: tt.fields.Repeated,
				Message:  tt.fields.Message,
				OneOf:    tt.fields.OneOf,
			}
			if got := template.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}
