package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/pkg/specs/types"
)

var gtypes = map[types.Type]graphql.Output{
	types.Double:   graphql.Float,
	types.Float:    graphql.Float,
	types.Uint64:   graphql.Int,
	types.Int64:    graphql.Int,
	types.Uint32:   graphql.Int,
	types.Int32:    graphql.Int,
	types.Fixed64:  graphql.Float,
	types.Fixed32:  graphql.Float,
	types.Bool:     graphql.Boolean,
	types.String:   graphql.String,
	types.Bytes:    graphql.String,
	types.Enum:     graphql.EnumValueType,
	types.Sfixed64: graphql.Float,
	types.Sfixed32: graphql.Float,
	types.Sint64:   graphql.Int,
	types.Sint32:   graphql.Int,
}
