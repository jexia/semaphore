package specs

import (
	"github.com/jexia/semaphore/pkg/specs/metadata"
)

// PropertyReference represents a mustach template reference
type PropertyReference struct {
	*metadata.Meta
	Resource string    `json:"resource,omitempty"`
	Path     string    `json:"path,omitempty"`
	Property *Property `json:"-"`
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
