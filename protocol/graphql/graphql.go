package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	log "github.com/sirupsen/logrus"
)

type req struct {
	Query string `json:"query"`
}

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) protocol.Listener {
	return &Listener{
		server: &http.Server{
			Addr: addr,
		},
	}
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

		req := req{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		result := graphql.Do(graphql.Params{
			Schema:        listener.schema,
			RequestString: req.Query,
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
	fields := graphql.Fields{}

	for _, endpoint := range endpoints {
		req := NewArgs(endpoint.Request.Property())
		res, err := NewObject(endpoint.Flow.Name, endpoint.Response.Property())
		if err != nil {
			return err
		}

		resolve := func(endpoint *protocol.Endpoint) graphql.FieldResolveFn {
			return func(p graphql.ResolveParams) (interface{}, error) {
				store := endpoint.Flow.NewStore()
				ctx := context.Background()

				store.StoreValues(specs.InputResource, "", p.Args)

				err = endpoint.Flow.Call(ctx, store)
				if err != nil {
					return nil, err
				}

				result, err := ResponseValue(endpoint.Response.Property(), store)
				if err != nil {
					return nil, err
				}

				return result, nil
			}
		}(endpoint)

		// TODO: set a option to set a custom name
		fields[endpoint.Flow.Name] = &graphql.Field{
			Args:    req,
			Type:    res,
			Resolve: resolve,
		}
	}

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   "query",
					Fields: fields,
				},
			),
		},
	)

	if err != nil {
		return err
	}

	listener.mutex.Lock()
	listener.schema = schema
	listener.mutex.Unlock()

	return nil
}

// Close closes the given listener
func (listener *Listener) Close() error {
	log.Info("Closing HTTP listener")
	return listener.server.Close()
}
