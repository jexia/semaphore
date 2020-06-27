package codec

import (
	"github.com/jexia/maestro/internal/codec/json"
	"github.com/jexia/maestro/internal/codec/proto"
)

// JSON constructs a new JSON message constructor
var JSON = json.NewConstructor

// Proto constructs a new Proto message constructor
var Proto = proto.NewConstructor
