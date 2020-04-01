package constructor

import (
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
)

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Ctx         instance.Context
	Definitions []specs.Resolver
	Codec       codec.Constructors
	Callers     transport.Callers
	Listeners   transport.Listeners
	Schemas     []schema.Resolver
	Schema      *schema.Store
	Functions   specs.CustomDefinedFunctions
}
