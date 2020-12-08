package transport

import "regexp"

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

	template, err = compileTemplate(template)
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

// compile template from `/<var_one>/<var_two>` to `/$var_one/$var_two`
// Workaround to make the template compatible with golang regexp variables
// (since we cannot use `$` (dollar) in HCL templates)
func compileTemplate(rawTemplate string) (string, error) {
	var (
		compiled []byte
		isOpened bool
	)

	for index, symbol := range rawTemplate {
		switch symbol {
		case '<':
			if isOpened {
				return "", ErrMalformedTemplate{
					Template: rawTemplate,
					Position: index,
					Cause:    "previous variable has not been closed",
				}
			}

			compiled, isOpened = append(compiled, '$'), true
		case '>':
			if !isOpened {
				return "", ErrMalformedTemplate{
					Template: rawTemplate,
					Position: index,
					Cause:    "variable has not been opened",
				}
			}

			isOpened = false
		default:
			compiled = append(compiled, byte(symbol))
		}
	}

	return string(compiled), nil
}
