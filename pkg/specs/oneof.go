package specs

// OneOf is a mixed type to let the schema validate values against exactly one of the properties.
// Example:
// OneOf{
//   "string":{Scalar: &Scalar{Type: types.String}},
//   "number":{Scalar: &Scalar{Type: types.Int32}},
//   "object":{Message: &Message{...}},
// }
// A given value must be one of these types: string, number or an object.
type OneOf map[string]*Property

func (oneOf OneOf) String() string { return dump(oneOf) }

// Clone OneOf.
func (oneOf OneOf) Clone() OneOf {
	return OneOf(Message(oneOf).Clone())
}

// Compare checks whether given OneOf mathches the provided one.
func (oneOf OneOf) Compare(expected OneOf) error {
	return Message(oneOf).Compare(Message(expected))
}
