package avros

import (
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

// AvroTypes is a lookup table for avro descriptor types
var AvroTypes = map[types.Type]string{
	types.Message: "record",
	types.Bool:    "boolean",
	types.Int32:   "int",
	types.Int64:   "int",
	types.Float:   "float",
	types.Double:  "double",
	types.Bytes:   "bytes",
	types.String:  "string",
	types.Enum:    "enum",
	types.Array:   "array",
}

// Types is a lookup table for avro descriptor types
var Types = map[string]types.Type{
	"record":  types.Message,
	"boolean": types.Bool,
	"int":     types.Int64,
	"float":   types.Float,
	"double":  types.Double,
	"bytes":   types.Bytes,
	"string":  types.String,
	"enum":    types.Enum,
	"array":   types.Array,
}
