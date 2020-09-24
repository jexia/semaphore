package hcl

import (
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

	var scalar = property.Scalar

	if scalar == nil {
		return errNonScalarType
	}

	switch value.Type() {
	case cty.String:
		scalar.Default = value.AsString()
		scalar.Type = types.String
	case cty.Number:
		var def int64
		gocty.FromCtyValue(value, &def)

		scalar.Default = def
		scalar.Type = types.Int64
	case cty.Bool:
		var def bool
		gocty.FromCtyValue(value, &def)

		scalar.Default = def
		scalar.Type = types.Bool
	default:
		return errUnknownPopertyType(value.Type().FriendlyName())
	}

	return nil
}
