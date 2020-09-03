package prettyerr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

const DefaultTextFormat = "{{ .Message }}\n{{ range $key, $value := .Details }}\t{{ $key }}: {{ JSON $value }}\n{{ end }}"

// TextFormatter represents Errors as a text.
// - stack is the Error stack
// - nodeTemplate is the template for a single stack element.
// Example: "({{.Code}}) {{.Message}}\n"
func TextFormatter(stack Errors, nodeTemplate string) (string, error) {
	funcs := template.FuncMap{
		"JSON": func(val interface{}) string {
			bytes, err := json.Marshal(val)
			if err != nil {
				return "<prettyerr: could not marshal the value>"
			}

			return string(bytes)
		},
	}

	tpl, err := template.New("node").Funcs(funcs).Parse(nodeTemplate)

	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	o := bytes.NewBufferString("")

	for i, pretty := range stack {
		err := tpl.Execute(o, pretty)
		if err != nil {
			return "", fmt.Errorf("failed to execute template for %d: %w", i, err)
		}
	}

	return o.String(), nil
}
