package hcl

// JoinPath joins the given flow paths
func JoinPath(values ...string) (result string) {
	for _, value := range values {
		if value == "" {
			continue
		}

		if len(result) > 0 {
			suffix := string(result[len(result)-1])
			if len(result) > 0 && (suffix != "." && suffix != ":") {
				result += "."
			}
		}

		result += value
	}

	if string(result[len(result)-1]) == "." {
		result = result[:len(result)-1]
	}

	return result
}
