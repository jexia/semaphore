package constructor

import (
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/definitions"
	"github.com/jexia/maestro/pkg/flow"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/transport"
)

// Option represents a constructor func which sets a given option
type Option func(*Options)

// NewOptions constructs a new options object
func NewOptions(ctx instance.Context) Options {
	return Options{
		Ctx:      ctx,
		Services: make([]definitions.ServicesResolver, 0),
		Flows:    make([]definitions.FlowsResolver, 0),
		Schemas:  make([]definitions.SchemaResolver, 0),
		Codec:    make(map[string]codec.Constructor),
	}
}

// Options represents all the available options
type Options struct {
	Ctx                   instance.Context
	Codec                 codec.Constructors
	Callers               transport.Callers
	Listeners             transport.Listeners
	Flows                 []definitions.FlowsResolver
	Endpoints             []definitions.EndpointsResolver
	Services              []definitions.ServicesResolver
	Schemas               []definitions.SchemaResolver
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
