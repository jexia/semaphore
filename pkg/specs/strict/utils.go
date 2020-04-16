package strict

import (
	"strings"
)

// GetService returns the service from the given endpoint
func GetService(path string) string {
	if strings.HasSuffix(path, ".") {
		path = path[:len(path)-1]
	}

	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		return parts[0]
	}

	return strings.Join(parts[:len(parts)-1], ".")
}

// GetMethod returns the method from the given endpoint
func GetMethod(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		return ""
	}

	return parts[len(parts)-1]
}
