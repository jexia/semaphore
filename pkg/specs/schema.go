package specs

// SchemaManifest represents a collection of messages which are used for type checks
type SchemaManifest struct {
	Properties map[string]*Property
}

// MergeSchemaManifest merges the incoming schema into the given schema manifest
func MergeSchemaManifest(left *SchemaManifest, incoming ...*SchemaManifest) {
	if left == nil {
		return
	}

	if left.Properties == nil {
		left.Properties = make(map[string]*Property, 0)
	}

	for _, manifest := range incoming {
		for key, val := range manifest.Properties {
			left.Properties[key] = val
		}
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
