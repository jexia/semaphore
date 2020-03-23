package specs

import (
	"context"

	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// SetDefaultValue sets the given value as default value inside the given property
func SetDefaultValue(ctx context.Context, property *Property, value cty.Value) {
	logger.FromCtx(ctx, logger.Core).WithField("path", property.Path).WithField("value", value).Debug("Set default value for property")

	switch value.Type() {
	case cty.String:
		property.Default = value.AsString()
		property.Type = types.TypeString
	case cty.Number:
		var def int64
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = types.TypeInt64
	case cty.Bool:
		var def bool
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = types.TypeBool
	}
}
