package specs

// NewSchemaManifest constructs a new schema manifest
func NewSchemaManifest() *SchemaManifest {
	return &SchemaManifest{
		Properties: make(map[string]*Property, 0),
	}
}

// SchemaManifest represents a collection of messages which are used for type checks
type SchemaManifest struct {
	Properties map[string]*Property
}

// Merge merges the incoming schema into the given schema manifest
func (schema *SchemaManifest) Merge(incoming *SchemaManifest) {
	if incoming == nil || schema == nil {
		return
	}

	if incoming.Properties == nil || schema.Properties == nil {
		return
	}

	for key, val := range incoming.Properties {
		schema.Properties[key] = val
	}
}

// GetProperty attempts to find a property with the given path inside the schema manifest
func (schema *SchemaManifest) GetProperty(path string) *Property {
	for key, prop := range schema.Properties {
		if key == path {
			return prop
		}
	}

	return nil
}
