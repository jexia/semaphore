package e2e

import (
	"io"
	"net/http"
	"testing"
)

// EchoHandler creates an HTTP handler that returns the request body as a response.
func EchoHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				t.Errorf("failed to close request body: %s", err)
			}
		}()

		if _, err := io.Copy(w, r.Body); err != nil {
			t.Errorf("failed to send the reply: %s", err)
		}
	}
}

// EchoRouter creates an HTTP router for testing.
func EchoRouter(t *testing.T) http.Handler {
	var mux = http.NewServeMux()

	mux.Handle("/echo", EchoHandler(t))
	// TODO: add more handlers

	return mux
}
