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
	return strings.Join(path[1:], ".")
}

// GetServiceProto attempts to find the service proto matching the given alias
func GetServiceProto(manifest *specs.Manifest, alias string) string {
	for _, service := range manifest.Services {
		if service.Alias == alias {
			return service.Proto
		}
	}

	return alias
}
