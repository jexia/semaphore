package template

import "strings"

// JoinPath joins the given flow paths
func JoinPath(values ...string) (result string) {
	for _, value := range values {
		if value == "" {
			continue
		}

		if len(result) > 0 {
			suffix := string(result[len(result)-1])
			if suffix != "." && suffix != ":" {
				result += "."
			}
		}

		result += value
	}

	if result == "" || result == "." {
		return result
	}

	if string(result[len(result)-1]) == "." {
		result = result[:len(result)-1]
	}

	if string(result[0]) == "." {
		result = result[1:]
	}

	return result
}

// SplitPath splits the given path into parts
func SplitPath(path string) []string {
	return strings.Split(path, PathDelimiter)
}

// ResourcePath constructs a new path using the given resource and path.
func ResourcePath(resource string, paths ...string) string {
	return resource + ReferenceDelimiter + JoinPath(paths...)
}
