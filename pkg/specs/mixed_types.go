package specs

// OneOf is a mixed type to let the schema validate values against exactly one of the templates.
// Example:
// OneOf{
//   {Scalar: &Scalar{Type: types.String}},
//   {Scalar: &Scalar{Type: types.Int32}},
//   {Message: &Message{...}},
// }
// A given value must be one of these types: string, int32 or the message.
type OneOf []*Template
