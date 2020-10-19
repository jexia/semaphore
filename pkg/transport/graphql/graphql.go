package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// Schema base
var (
	QueryObject    = "query"
	MutationObject = "mutation"
)

type req struct {
	Query string `json:"query"`
}

// NewListener constructs a new listener for the given addr
func NewListener(addr string, opts specs.Options) transport.NewListener {
	return func(parent *broker.Context) transport.Listener {
		module := broker.WithModule(parent, "listener", "graphql")
		ctx := logger.WithLogger(module)

		return &Listener{
			ctx: ctx,
			server: &http.Server{
				Addr: addr,
			},
		}
	}
}

// Listener represents a GraphQL listener
type Listener struct {
	ctx    *broker.Context
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

	logger.Info(listener.ctx, "serving GraphQL listener", zap.String("addr", listener.server.Addr))

	err := listener.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Handle parses the given endpoints and constructs route handlers
func (listener *Listener) Handle(ctx *broker.Context, endpoints []*transport.Endpoint, constructors map[string]codec.Constructor) error {
	if len(endpoints) == 0 {
		listener.mutex.Lock()
		listener.schema = graphql.Schema{}
		listener.mutex.Unlock()
		return nil
	}

	objects := NewObjects()
	fields := map[string]graphql.Fields{
		QueryObject:    {},
		MutationObject: {},
	}

	for _, endpoint := range endpoints {
		options, err := ParseEndpointOptions(endpoint)
		if err != nil {
			return err
		}

		err = endpoint.NewCodec(ctx, nil, nil)
		if err != nil {
			return err
		}

		resolve := func(endpoint *transport.Endpoint) graphql.FieldResolveFn {
			return func(p graphql.ResolveParams) (interface{}, error) {
				store := endpoint.Flow.NewStore()
				ctx := context.Background()

				store.StoreValues(template.InputResource, "", p.Args)

				err = endpoint.Flow.Do(ctx, store)
				if err != nil {
					object := endpoint.Errs.Get(transport.Unwrap(err))
					if object == nil {
						logger.Error(listener.ctx, "unable to lookup error manager", zap.Error(err))
						return nil, err
					}

					message := object.ResolveMessage(store)
					return nil, errors.New(message)
				}

				if endpoint.Response == nil || endpoint.Response.Definition == nil {
					return make(map[string]interface{}), nil
				}

				result, err := ResponseObject(endpoint.Response.Definition.Property, store)
				if err != nil {
					return nil, err
				}

				return result, nil
			}
		}(endpoint)

		res, err := NewSchemaObject(objects, options.Name, endpoint.Response)
		if err != nil {
			return err
		}

		path := options.Path
		field := &graphql.Field{
			Args:    graphql.FieldConfigArgument{},
			Resolve: resolve,
			Type:    res,
		}

		if endpoint.Request != nil {
			req, err := NewArgs(endpoint.Request.Definition)
			if err != nil {
				return err
			}

			field.Args = req

			if endpoint.Request.Definition != nil && endpoint.Request.Definition.Property != nil {
				field.Description = endpoint.Request.Definition.Property.Description
			}
		}

		if options.Base == QueryObject && field.Type == nil {
			options.Base = MutationObject
		}

		err = SetField(path, fields[options.Base], field)
		if err != nil {
			return err
		}
	}

	config := graphql.SchemaConfig{}

	if len(fields[MutationObject]) > 0 {
		config.Mutation = graphql.NewObject(
			graphql.ObjectConfig{
				Name:   MutationObject,
				Fields: fields[MutationObject],
			},
		)
	}

	if len(fields[QueryObject]) > 0 {
		config.Query = graphql.NewObject(
			graphql.ObjectConfig{
				Name:   QueryObject,
				Fields: fields[QueryObject],
			},
		)
	}

	schema, err := graphql.NewSchema(config)
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
	logger.Info(listener.ctx, "closing GraphQL listener")
	return listener.server.Close()
}
