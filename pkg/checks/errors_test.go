package checks

import (
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

func TestErrFlowDuplicate_Error(t *testing.T) {
	type fields struct {
		Flow string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return formatted error",
			fields{Flow: "WalkTheDog"},
			"duplicate flow 'WalkTheDog'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrFlowDuplicate{
				Flow: tt.fields.Flow,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrFlowDuplicate_Prettify(t *testing.T) {
	type fields struct {
		Flow string
	}
	tests := []struct {
		name   string
		fields fields
		want   prettyerr.Error
	}{
		{
			"return pretty error",
			fields{Flow: "WalkTheDog"},
			prettyerr.Error{
				Message: "duplicate flow 'WalkTheDog'",
				Details: map[string]interface{}{
					"flow": "WalkTheDog",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrFlowDuplicate{
				Flow: tt.fields.Flow,
			}
			if got := e.Prettify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prettify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrReservedKeyword_Error(t *testing.T) {
	type fields struct {
		Flow string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return formatted error",
			fields{Flow: "input"},
			"flow with the name 'input' is a reserved keyword",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrReservedKeyword{
				Flow: tt.fields.Flow,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrReservedKeyword_Prettify(t *testing.T) {
	type fields struct {
		Flow string
	}
	tests := []struct {
		name   string
		fields fields
		want   prettyerr.Error
	}{
		{
			"return pretty error",
			fields{Flow: "input"},
			prettyerr.Error{
				Message: "flow with the name 'input' is a reserved keyword",
				Details: map[string]interface{}{"flow": "input"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrReservedKeyword{
				Flow: tt.fields.Flow,
			}
			if got := e.Prettify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prettify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrResourceDuplicate_Error(t *testing.T) {
	type fields struct {
		Flow     string
		Resource string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return formatted error",
			fields{Flow: "WalkTheDog", Resource: "Alva"},
			"duplicate resource 'Alva' in flow 'WalkTheDog'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrResourceDuplicate{
				Flow:     tt.fields.Flow,
				Resource: tt.fields.Resource,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrResourceDuplicate_Prettify(t *testing.T) {
	type fields struct {
		Flow     string
		Resource string
	}
	tests := []struct {
		name   string
		fields fields
		want   prettyerr.Error
	}{
		{
			"return pretty error",
			fields{Flow: "WalkTheDog", Resource: "Alva"},
			prettyerr.Error{
				Message: "duplicate resource 'Alva' in flow 'WalkTheDog'",
				Details: map[string]interface{}{"flow": "WalkTheDog", "resource": "Alva"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrResourceDuplicate{
				Flow:     tt.fields.Flow,
				Resource: tt.fields.Resource,
			}
			if got := e.Prettify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prettify() = %v, want %v", got, tt.want)
			}
		})
	}
}
