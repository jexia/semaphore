package labels

import "strings"

var labels = map[Label]string{
	Optional: "optional",
	Required: "required",
}

var keys = map[string]Label{
	"optional": Optional,
	"required": Required,
}

const delimiter = ","

// Label represents a value label
type Label int

func (label Label) String() string {
	msg := strings.Builder{}

	for l, v := range labels {
		if label&l == 0 {
			if msg.Len() > 0 {
				msg.WriteString(delimiter + " ")
			}

			msg.WriteString(v)
		}
	}

	return msg.String()
}

// Has checks whether the given key is available inside the given label
func (label Label) Has(key Label) bool {
	return label&key == 0
}

const (
	// Optional representing a optional field
	Optional Label = 1 << iota
	// Required representing a required field
	Required
)

// Parse parses the given label string into a label
func Parse(label string) (result Label) {
	sliced := strings.Split(label, delimiter)
	for _, raw := range sliced {
		key := strings.TrimSpace(raw)
		result = result | keys[key]
	}

	return result
}
