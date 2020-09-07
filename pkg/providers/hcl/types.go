package hcl

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"go.uber.org/zap"
)

// SetDefaultValue sets the given value as default value inside the given property
func SetDefaultValue(ctx *broker.Context, property *specs.Property, value cty.Value) error {
	logger.Debug(ctx, "set default value for property", zap.String("path", property.Path), zap.Any("value", value))

	switch value.Type() {
	case cty.String:
		property.Default = value.AsString()
		property.Type = types.String
	case cty.Number:
		var def int64
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = types.Int64
	case cty.Bool:
		var def bool
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = types.Bool
	default:
		return fmt.Errorf("unknown property type: %T", value.Type())
	}

	return nil
}
