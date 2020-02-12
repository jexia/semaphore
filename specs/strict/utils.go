package strict

import (
	"strings"

	"github.com/jexia/maestro/specs"
)

// GetService returns the service from the given endpoint
func GetService(endpoint string) string {
	path := strings.Split(endpoint, ".")
	return path[0]
}

// GetMethod returns the method from the given endpoint
func GetMethod(endpoint string) string {
	path := strings.Split(endpoint, ".")
	if len(path) <= 1 {
		return ""
	}

	return strings.Join(path[1:], ".")
}

// GetSchemaService attempts to find the schema service matching the given alias
func GetSchemaService(manifest *specs.Manifest, alias string) string {
	for _, service := range manifest.Services {
		if service.Alias == alias {
			return service.Schema
		}
	}

	return alias
}
