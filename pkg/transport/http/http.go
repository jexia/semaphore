package http

import (
	"regexp"
)

// ReferenceLookup is executed to lookup references within a given endpoint.
// A reference starts with a colon and could contain characters, numbers, underscores and hyphens and dots to define nested properties
var ReferenceLookup = regexp.MustCompile(`(?m):[a-zA-Z\d\^\&\%\$@\_\-\.]+`)
