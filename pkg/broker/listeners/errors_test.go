package listeners

import (
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

func TestErrNoListener_Error(t *testing.T) {
	type fields struct {
		Listener string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return formatted error",
			fields{Listener: "dogs"},
			"unknown listener 'dogs'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrNoListener{
				Listener: tt.fields.Listener,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoListener_Prettify(t *testing.T) {
	type fields struct {
		Listener string
	}
	tests := []struct {
		name   string
		fields fields
		want   prettyerr.Error
	}{
		{
			"return pretty error",
			fields{Listener: "dogs"},
			prettyerr.Error{
				Message: "unknown listener 'dogs'",
				Details: map[string]interface{}{"listener": "dogs"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrNoListener{
				Listener: tt.fields.Listener,
			}
			if got := e.Prettify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prettify() = %v, want %v", got, tt.want)
			}
		})
	}
}
