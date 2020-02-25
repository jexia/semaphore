package http

import (
	"net/http"

	"github.com/jexia/maestro/protocol"
)

// CopyHeader copies the given protocol header into a HTTP header
func CopyHeader(header protocol.Header) http.Header {
	result := http.Header{}
	for key, val := range header {
		result.Set(key, val)
	}

	return result
}
