package http

import (
	"regexp"
)

// ReferenceLookup is executed to lookup references within a given endpoint
var ReferenceLookup = regexp.MustCompile(`(?m):\w+`)
