package providers

import (
	"errors"
	"testing"
)

func TestErrUndefinedObject(t *testing.T) {
	type fields struct {
		Schema string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Schema: "get"},
			"object 'get', is unavailable inside the schema collection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedObject{
				Schema: "get",
			}

			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedService(t *testing.T) {
	type fields struct {
		Service string
		Flow    string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Service: "getsources", Flow: "add"},
			"undefined service 'getsources' in flow 'add'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedService{
				Service: "getsources",
				Flow:    "add",
			}

			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedMethod(t *testing.T) {
	type fields struct {
		Method string
		Flow   string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Method: "get", Flow: "add"},
			"undefined method 'get' in flow 'add'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedMethod{
				Method: "get",
				Flow:   "add",
			}

			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedOutput(t *testing.T) {
	type fields struct {
		Output string
		Flow   string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Output: "json", Flow: "add"},
			"undefined method output property 'json' in flow 'add'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedOutput{
				Output: "json",
				Flow:   "add",
			}

			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedProperty(t *testing.T) {
	type fields struct {
		Property string
		Flow     string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Property: "getdata", Flow: "add"},
			"undefined schema nested message property 'getdata' in flow 'add': something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedProperty{
				Property: "getdata",
				Flow:     "add",
				Inner:    errors.New("something went wrong"),
			}

			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
