package types

// Label represents a value label
type Label string

// Spec labels
const (
	LabelOptional Label = "optional"
	LabelRequired Label = "required"
	LabelRepeated Label = "repeated"
)
