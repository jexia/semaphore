package labels

// Label represents a value label.
type Label string

const (
	// Optional representing a optional field.
	Optional Label = "optional"

	// Required representing a required field.
	Required Label = "required"
)
