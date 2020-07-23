package metadata

type empty struct{}

func (empty) Value(key interface{}) interface{} {
	return nil
}

// Empty returns a empty meta object
func Empty() Meta {
	return empty{}
}

// Meta represents a metadata store capable of holding values that
// extend a given struct through a context Value like approach.
type Meta interface {
	Value(key interface{}) interface{}
}

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent Meta, key, value interface{}) Meta {
	if parent == nil {
		parent = Empty()
	}

	return &valueMeta{parent, key, value}
}

type valueMeta struct {
	Meta
	key interface{}
	val interface{}
}

func (c *valueMeta) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}

	return c.Meta.Value(key)
}

// WithValues combines the two given meta objects into one
func WithValues(parent Meta, child Meta) Meta {
	if parent == nil && child == nil {
		return Empty()
	}

	if child == nil {
		return parent
	}

	if parent == nil {
		return child
	}

	return &valuesMeta{parent, child}
}

type valuesMeta struct {
	parent Meta
	child  Meta
}

func (c *valuesMeta) Value(key interface{}) interface{} {
	v := c.parent.Value(key)
	if v != nil {
		return v
	}

	return c.child.Value(key)
}
