package graphql

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
)

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) (protocol.Listener, error) {
	return &Listener{
		server: &http.Server{
			Addr: addr,
		},
	}, nil
}

// Listener represents a GraphQL listener
type Listener struct {
	schema graphql.Schema
	mutex  sync.RWMutex
	server *http.Server
}

// Name returns the name of the given listener
func (listener *Listener) Name() string {
	return "graphql"
}

// Serve opens the GraphQL listener and calls the given handler function on reach request
func (listener *Listener) Serve() error {
	listener.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listener.mutex.RLock()
		defer listener.mutex.RUnlock()

		query := r.URL.Query().Get("query")
		result := graphql.Do(graphql.Params{
			Schema:        listener.schema,
			RequestString: query,
		})

		json.NewEncoder(w).Encode(result)
	})

	err := listener.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(endpoints []*protocol.Endpoint) error {

	log.Println(endpoints)

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	return nil
}
