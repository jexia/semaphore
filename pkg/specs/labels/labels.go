package labels

import "strings"

var labels = map[Label]string{
	Optional: "optional",
	Required: "required",
}

// Label represents a value label
type Label int

func (label Label) String() string {
	msg := strings.Builder{}

	for l, v := range labels {
		if label&l == 0 {
			if msg.Len() > 0 {
				msg.WriteString(", ")
			}

			msg.WriteString(v)
		}
	}

	return msg.String()
}

func (label Label) Has(key Label) bool {
	return label&key == 0
}

// Spec labels
const (
	Optional Label = 1 << iota
	Required
)

func Compatible(pattern, value Label) bool {
	return true
}
