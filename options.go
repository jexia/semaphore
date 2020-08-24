package semaphore

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Option represents a constructor func which sets a given option
type Option func(*broker.Context, *Options)

// Options represents all the available options
type Options struct {
	Codec                 codec.Constructors
	Callers               transport.Callers
	FlowResolvers         providers.FlowsResolvers
	Middleware            []Middleware
	BeforeConstructor     BeforeConstructor
	BeforeManagerDo       flow.BeforeManager
	BeforeManagerRollback flow.BeforeManager
	AfterManagerDo        flow.AfterManager
	AfterManagerRollback  flow.AfterManager
	BeforeNodeDo          flow.BeforeNode
	BeforeNodeRollback    flow.BeforeNode
	AfterNodeDo           flow.AfterNode
	AfterNodeRollback     flow.AfterNode
	AfterFlowConstruction AfterFlowConstruction
	Functions             functions.Custom
}

// Middleware is called once the options have been initialised
type Middleware func(*broker.Context) ([]Option, error)

// BeforeConstructor is called before the specifications is constructored
type BeforeConstructor func(*broker.Context, functions.Collection, Options) error

// BeforeConstructorHandler wraps the before constructed function to allow middleware to be chained
type BeforeConstructorHandler func(BeforeConstructor) BeforeConstructor

// AfterFlowConstruction is called before the construction of a flow manager
type AfterFlowConstruction func(*broker.Context, specs.FlowInterface, *flow.Manager) error

// AfterFlowConstructionHandler wraps the before flow construction function to allow middleware to be chained
type AfterFlowConstructionHandler func(AfterFlowConstruction) AfterFlowConstruction

// NewOptions constructs a Options object from the given Option constructors
func NewOptions(ctx *broker.Context, options ...Option) (Options, error) {
	result := Options{
		FlowResolvers: make([]providers.FlowsResolver, 0),
		Codec:         make(map[string]codec.Constructor),
	}

	if options == nil {
		return result, nil
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(ctx, &result)
	}

	for _, middleware := range result.Middleware {
		options, err := middleware(ctx)
		if err != nil {
			return result, err
		}

		for _, option := range options {
			option(ctx, &result)
		}
	}

	return result, nil
}

// NewCollection constructs a new options collection
func NewCollection(options ...Option) []Option {
	return options
}

// WithFlows appends the given flows resolver to the available flow resolvers
func WithFlows(definition providers.FlowsResolver) Option {
	return func(ctx *broker.Context, options *Options) {
		options.FlowResolvers = append(options.FlowResolvers, definition)
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(codec codec.Constructor) Option {
	return func(ctx *broker.Context, options *Options) {
		options.Codec[codec.Name()] = codec
	}
}

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller transport.NewCaller) Option {
	return func(ctx *broker.Context, options *Options) {
		options.Callers = append(options.Callers, caller(ctx))
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(custom functions.Custom) Option {
	return func(ctx *broker.Context, options *Options) {
		if options.Functions == nil {
			options.Functions = functions.Custom{}
		}

		for key, fn := range custom {
			options.Functions[key] = fn
		}
	}
}

// WithLogLevel sets the log level for the given module
func WithLogLevel(pattern string, value string) Option {
	return func(ctx *broker.Context, options *Options) {
		level := zapcore.InfoLevel
		err := level.UnmarshalText([]byte(value))
		if err != nil {
			logger.Error(ctx, "unable to unmarshal log level", zap.String("level", value))
		}

		logger.SetLevel(ctx, pattern, level)
	}
}
