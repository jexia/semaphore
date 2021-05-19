package specs

import (
	"encoding/json"
	"strings"
)

func dump(value interface{}) string {
	var (
		buff    strings.Builder
		encoder = json.NewEncoder(&buff)
	)

	encoder.SetIndent("", "  ")
	encoder.Encode(value)

	return buff.String()
}
