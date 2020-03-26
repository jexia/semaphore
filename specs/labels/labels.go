package labels

// Label represents a value label
type Label string

// Spec labels
const (
	Optional Label = "optional"
	Required Label = "required"
	Repeated Label = "repeated"
)
