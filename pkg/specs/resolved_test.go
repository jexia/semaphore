package specs

import "testing"

func TestResolvedProperty(t *testing.T) {
	var (
		resolved = NewResolvedProperty()

		propertyA = &Property{
			Path:     "property.A",
			Template: new(Template),
		}

		propertyB = &Property{
			Template: &Template{
				Identifier: "B",
			},
		}

		propertyC = new(Property)
	)

	if resolved.Resolved(propertyA) {
		t.Errorf("propertyA should not be resolved")
	}

	if resolved.Resolved(propertyB) {
		t.Errorf("propertyB should not be resolved")
	}

	resolved.Resolve(propertyA)
	resolved.Resolve(propertyB)
	resolved.Resolve(propertyC)

	if !resolved.Resolved(propertyA) {
		t.Errorf("propertyA should be resolved")
	}

	if !resolved.Resolved(propertyB) {
		t.Errorf("propertyB should be resolved")
	}

	if resolved.Resolved(propertyC) {
		t.Errorf("propertyB should never be resolved")
	}
}
