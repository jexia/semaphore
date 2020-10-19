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

	defaultTpl, _ := template.New("node").Funcs(funcs).Parse(DefaultTextFormat)
	tpl, err := template.New("node").Funcs(funcs).Parse(nodeTemplate)

	if err != nil {
		err = fmt.Errorf("failed to parse template: %w", err)
		// Append failed to parse template error to stack
		stack = append(Errors{Error{
			Original: err,
			Message:  err.Error(),
			Details:  nil,
			Code:     GenericErrorCode,
		}}, stack...)

		// Fallback to default template
		tpl = defaultTpl
	}

	o := bytes.NewBufferString("")

	for i, pretty := range stack {
		err := tpl.Execute(o, pretty)
		if err != nil {
			err = fmt.Errorf("failed to execute template for %d: %w", i, err)

			// Append failed to execute template error to string using default template
			_ = defaultTpl.Execute(o, Error{
				Original: err,
				Message:  err.Error(),
				Details:  nil,
				Code:     GenericErrorCode,
			})
			_ = defaultTpl.Execute(o, pretty)
		}
	}

	return o.String(), nil
}
