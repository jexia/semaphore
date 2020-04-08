package hcl

import (
	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// SetDefaultValue sets the given value as default value inside the given property
func SetDefaultValue(ctx instance.Context, property *specs.Property, value cty.Value) {
	ctx.Logger(logger.Core).WithField("path", property.Path).WithField("value", value).Debug("Set default value for property")

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
	}
}
