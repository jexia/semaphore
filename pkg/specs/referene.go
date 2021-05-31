package specs

import (
	"strings"

	"github.com/jexia/semaphore/v2/pkg/specs/metadata"
)

// ReferenceDelimiter represents the value resource reference delimiter.
const ReferenceDelimiter = ":"

// PropertyReference represents a mustach template reference
type PropertyReference struct {
	*metadata.Meta
	Resource string    `json:"resource,omitempty"`
	Path     string    `json:"path,omitempty"`
	Property *Property `json:"-"`
}

// ParsePropertyReference parses the given value to a property reference.
func ParsePropertyReference(value string) *PropertyReference {
	var prop string

	rv := strings.Split(value, ReferenceDelimiter)
	resource := rv[0]
	resources := SplitPath(resource)

	if len(resources) > 1 {
		prop = JoinPath(resources[1:]...)
	}

	reference := &PropertyReference{
		Resource: resource,
	}

	if len(rv) == 1 {
		return reference
	}

	path := rv[1]

	if prop == HeaderResource {
		path = strings.ToLower(path)
	}

	reference.Path = path
	return reference
}

// Clone clones the given property reference
func (reference *PropertyReference) Clone() *PropertyReference {
	if reference == nil {
		return nil
	}

	return &PropertyReference{
		Meta:     reference.Meta,
		Resource: reference.Resource,
		Path:     reference.Path,
		Property: nil,
	}
}

func (reference *PropertyReference) String() string {
	if reference == nil {
		return ""
	}

	return reference.Resource + ":" + reference.Path
}
