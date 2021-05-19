package types

import "github.com/jexia/semaphore/v2/pkg/specs/types"

// Type represents a openapi3 type
type Type string

// Available openapi3 types
const (
	String  Type = "string"
	Number  Type = "number"
	Integer Type = "integer"
	Boolean Type = "boolean"
	Array   Type = "array"
	Object  Type = "object"
)

var s2oTypes = map[types.Type]Type{
	types.Double:   Integer,
	types.Int64:    Integer,
	types.Uint64:   Integer,
	types.Int32:    Integer,
	types.Uint32:   Integer,
	types.Fixed32:  Integer,
	types.Fixed64:  Integer,
	types.Float:    Integer,
	types.String:   String,
	types.Enum:     String,
	types.Bool:     Boolean,
	types.Bytes:    String,
	types.Sfixed32: Integer,
	types.Sfixed64: Integer,
	types.Sint32:   Integer,
	types.Sint64:   Integer,
}

var o2sTypes = make(map[Type]types.Type, len(s2oTypes))

func init() {
	// ensure that both types are set
	for t, o := range s2oTypes {
		o2sTypes[o] = t
	}
}

// Open returns the representing OpenApi 3 type for the given specification type
func Open(tp types.Type) Type {
	return s2oTypes[tp]
}
