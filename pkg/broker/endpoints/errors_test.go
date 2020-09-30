package endpoints

import (
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

func TestErrNoServiceForMethod_Error(t *testing.T) {
	type fields struct {
		Method string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Method: "get"},
			"failed to find service for 'get'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrNoServiceForMethod{
				Method: tt.fields.Method,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoServiceForMethod_Prettify(t *testing.T) {
	type fields struct {
		Method string
	}
	tests := []struct {
		name   string
		fields fields
		want   prettyerr.Error
	}{
		{
			"return pretty error",
			fields{Method: "get"},
			prettyerr.Error{Message: "failed to find service for 'get'", Details: map[string]interface{}{"method": "get"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrNoServiceForMethod{
				Method: tt.fields.Method,
			}
			if got := e.Prettify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prettify() = %v, want %v", got, tt.want)
			}
		})
	}
}
