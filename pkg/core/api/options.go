package api

import (
	"github.com/jexia/maestro/internal/codec"
	"github.com/jexia/maestro/internal/flow"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/providers"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

// Collection represents a collection of manifests and functions
type Collection struct {
	Flows     *specs.FlowsManifest
	Endpoints *specs.EndpointsManifest
	Services  *specs.ServicesManifest
	Schema    *specs.SchemaManifest
}

// Option represents a constructor func which sets a given option
type Option func(*Options)

// NewOptions constructs a new options object
func NewOptions(ctx instance.Context) Options {
	return Options{
		Ctx:      ctx,
		Services: make([]providers.ServicesResolver, 0),
		Flows:    make([]providers.FlowsResolver, 0),
		Schemas:  make([]providers.SchemaResolver, 0),
		Codec:    make(map[string]codec.Constructor),
	}
}

// Options represents all the available options
type Options struct {
	Ctx                   instance.Context
	Codec                 codec.Constructors
	Callers               transport.Callers
	Listeners             transport.Listeners
	Flows                 []providers.FlowsResolver
	Endpoints             []providers.EndpointsResolver
	Services              []providers.ServicesResolver
	Schemas               []providers.SchemaResolver
	Middleware            []Middleware
	AfterConstructor      AfterConstructor
	BeforeManagerDo       flow.BeforeManager
	BeforeManagerRollback flow.BeforeManager
	AfterManagerDo        flow.AfterManager
	AfterManagerRollback  flow.AfterManager
	BeforeNodeDo          flow.BeforeNode
	BeforeNodeRollback    flow.BeforeNode
	AfterNodeDo           flow.AfterNode
	AfterNodeRollback     flow.AfterNode
	Functions             functions.Custom
}

// Middleware is called once the options have been initialised
type Middleware func(instance.Context) ([]Option, error)

// AfterConstructor is called after the specifications is constructored
type AfterConstructor func(instance.Context, *Collection) error

// AfterConstructorHandler wraps the after constructed function to allow middleware to be chained
type AfterConstructorHandler func(AfterConstructor) AfterConstructor
