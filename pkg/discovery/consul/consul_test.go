package consul

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolver(t *testing.T) {
	t.Parallel()

	// Run a mock HTTPserver to allow the watcher to connect
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	t.Cleanup(srv.Close)

	t.Run("should return a watcher", func(t *testing.T) {
		t.Parallel()
		consul := New()

		t.Cleanup(func() {
			for _, watcher := range consul.watchers {
				watcher.Stop()
			}
		})

		watcher, err := consul.Resolver(srv.URL)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if watcher == nil {
			t.Fatal("expected to get a new watcher")
		}

		if len(consul.watchers) != 1 {
			t.Fatal("expected to have exactly one watcher in the consul manager")
		}
	})

	t.Run("should return a cached watcher", func(t *testing.T) {
		t.Parallel()
		consul := New()

		t.Cleanup(func() {
			for _, watcher := range consul.watchers {
				watcher.Stop()
			}
		})

		watcher, err := consul.Resolver(srv.URL)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if watcher == nil {
			t.Fatal("expected to get a new watcher")
		}

		cachedWatcher, err := consul.Resolver(srv.URL)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if cachedWatcher != watcher {
			t.Fatal("expected to get the same watcher from cache")
		}

		if len(consul.watchers) != 1 {
			t.Fatal("expected to have exactly one watcher in the consul manager")
		}
	})
}
