package specs

// SchemaManifest represents a collection of messages which are used for type checks
type SchemaManifest struct {
	Properties Objects `json:"properties"`
}

// Append merges the incoming schema into the given schema manifest
func (manifest *SchemaManifest) Append(incoming ...*SchemaManifest) {
	if manifest == nil {
		return
	}

	if manifest.Properties == nil {
		manifest.Properties = make(map[string]*Property, 0)
	}

	for _, right := range incoming {
		for key, val := range right.Properties {
			manifest.Properties[key] = val
		}
	}
}

// GetProperty attempts to find a property with the given path inside the schema manifest
func (manifest *SchemaManifest) GetProperty(path string) *Property {
	for key, prop := range manifest.Properties {
		if key == path {
			return prop
		}
	}

	return nil
}
