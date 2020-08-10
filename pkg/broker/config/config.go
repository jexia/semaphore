package config

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

// Constructor represents a specs constructor that could be used to construct specifications
type Constructor func(ctx *broker.Context, mem functions.Collection, options Options) (specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas, error)

// Option represents a constructor func which sets a given option
type Option func(*Options)

// NewOptions constructs a new options object
func NewOptions(ctx *broker.Context) Options {
	return Options{
		Ctx:              ctx,
		ServiceResolvers: make([]providers.ServicesResolver, 0),
		FlowResolvers:    make([]providers.FlowsResolver, 0),
		SchemaResolvers:  make([]providers.SchemaResolver, 0),
		Codec:            make(map[string]codec.Constructor),
	}
}

// Options represents all the available options
type Options struct {
	Ctx                   *broker.Context
	Constructor           Constructor
	Codec                 codec.Constructors
	Callers               transport.Callers
	Listeners             transport.Listeners
	FlowResolvers         providers.FlowsResolvers
	EndpointResolvers     providers.EndpointResolvers
	ServiceResolvers      providers.ServiceResolvers
	SchemaResolvers       providers.SchemaResolvers
	Middleware            []Middleware
	BeforeConstructor     BeforeConstructor
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
type Middleware func(*broker.Context) ([]Option, error)

// AfterConstructor is called after the specifications is constructored
type AfterConstructor func(*broker.Context, specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas) error

// AfterConstructorHandler wraps the after constructed function to allow middleware to be chained
type AfterConstructorHandler func(AfterConstructor) AfterConstructor

// BeforeConstructor is called before the specifications is constructored
type BeforeConstructor func(*broker.Context, functions.Collection, Options) error

// BeforeConstructorHandler wraps the before constructed function to allow middleware to be chained
type BeforeConstructorHandler func(BeforeConstructor) BeforeConstructor
