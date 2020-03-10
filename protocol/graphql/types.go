package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/specs/types"
)

var gtypes = map[types.Type]graphql.Output{
	types.TypeDouble:   graphql.Float,
	types.TypeFloat:    graphql.Float,
	types.TypeUint64:   graphql.Int,
	types.TypeInt64:    graphql.Int,
	types.TypeUint32:   graphql.Int,
	types.TypeInt32:    graphql.Int,
	types.TypeFixed64:  graphql.Float,
	types.TypeFixed32:  graphql.Float,
	types.TypeBool:     graphql.Boolean,
	types.TypeString:   graphql.String,
	types.TypeBytes:    graphql.String,
	types.TypeSfixed64: graphql.Float,
	types.TypeSfixed32: graphql.Float,
	types.TypeSint64:   graphql.Int,
	types.TypeSint32:   graphql.Int,
}
