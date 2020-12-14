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

// SetScalar sets the given value as default value inside the given property
func SetScalar(ctx *broker.Context, tmpl *specs.Template, value cty.Value) error {
	logger.Debug(ctx, "set default value for property", zap.Any("value", value))

	scalar := &specs.Scalar{}

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
		return ErrUnkownPropertyType(value.Type().FriendlyName())
	}

	tmpl.Scalar = scalar
	return nil
}
