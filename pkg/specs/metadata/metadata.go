package metadata

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent *Meta, key, value interface{}) *Meta {
	return &Meta{parent, key, value}
}

// Meta represents a metadata store capable of holding values that
// extend a given struct through a context Value like approach.
type Meta struct {
	parent *Meta
	key    interface{}
	val    interface{}
}

// Value compares the given key with the value key inside
// the given meta object. If the keys match is the value
// returned. If the keys do not match is the parent object called.
func (c *Meta) Value(key interface{}) interface{} {
	if c == nil {
		return nil
	}

	if c.key == key {
		return c.val
	}

	return c.parent.Value(key)
}
