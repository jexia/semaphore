package strict

import (
	"strings"
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
