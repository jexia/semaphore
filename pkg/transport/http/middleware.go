package http

import (
	"net/http"

	"github.com/rs/cors"
)

// OptionsHandler creates a global handler for OPTIONS requests.
func OptionsHandler(origins, headers, methods []string) http.Handler {
	var allowOrigin func(origin string) bool

	if len(origins) > 0 {
		allowOrigin = func(origin string) bool {
			for _, curr := range origins {
				if curr == origin {
					return true
				}
			}

			return false
		}
	}

	mw := cors.New(cors.Options{
		AllowOriginFunc: allowOrigin,
		AllowedHeaders:  headers,
		AllowedMethods:  methods,
		Debug:           true,
	})

	return mw.Handler(
		http.HandlerFunc(
			func(http.ResponseWriter, *http.Request) {},
		),
	)
}
