package transport

import (
	"regexp"
)

// Rewrite describes a function which is applied to modify the request URL.
// Returns modified URL and true or source URL and false if does not match.
type Rewrite func(string) (string, bool)

// NewRewrite prepares a rewrite function to be called by the handler
// before forwarding the request to the proxy.
func NewRewrite(source, template string) (Rewrite, error) {
	expression, err := regexp.Compile(source)
	if err != nil {
		return nil, err
	}

	return func(sourceURL string) (string, bool) {
		if !expression.MatchString(sourceURL) {
			return sourceURL, false
		}

		var targetURL = make([]byte, 0)

		for _, submatches := range expression.FindAllStringSubmatchIndex(sourceURL, -1) {
			targetURL = expression.ExpandString(targetURL, template, sourceURL, submatches)
		}

		return string(targetURL), true
	}, nil
}
