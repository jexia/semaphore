package specs

import "testing"

func TestServicesAppend(t *testing.T) {
	services := ServiceList{}

	services.Append(ServiceList{&Service{}, &Service{}})

	if len(services) != 2 {
		t.Fatalf("unexpected length %+v, expected 2", len(services))
	}
}

func TestServicesAppendNilValue(t *testing.T) {
	var services ServiceList
	services.Append(nil)
}

func TestServicesGet(t *testing.T) {
	services := ServiceList{&Service{FullyQualifiedName: "first"}, &Service{FullyQualifiedName: "second"}}

	result := services.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestServicesGetUnknown(t *testing.T) {
	services := ServiceList{&Service{FullyQualifiedName: "first"}}

	result := services.Get("unknown")
	if result != nil {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestServiceGetMethod(t *testing.T) {
	service := &Service{
		Methods: []*Method{
			{
				Name: "first",
			},
			{
				Name: "second",
			},
		},
	}

	result := service.GetMethod("second")
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestServiceGetUnknownMethod(t *testing.T) {
	service := &Service{
		Methods: []*Method{
			{
				Name: "first",
			},
		},
	}

	result := service.GetMethod("unknown")
	if result != nil {
		t.Fatalf("unexpected result %+v", result)
	}
}
