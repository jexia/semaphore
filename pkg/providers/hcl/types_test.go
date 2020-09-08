package hcl

import (
	"errors"
	"math/big"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/zclconf/go-cty/cty"
)

func TestSetDefaultValue(t *testing.T) {
	type expected struct {
		defaultValue interface{}
		dataType     types.Type
		error        error
	}

	tests := map[cty.Value]expected{
		cty.StringVal("default"): {
			defaultValue: "default",
			dataType:     types.String,
		},
		cty.NumberVal(big.NewFloat(10)): {
			defaultValue: int64(10),
			dataType:     types.Int64,
		},
		cty.BoolVal(true): {
			defaultValue: true,
			dataType:     types.Bool,
		},
		cty.DynamicVal: {
			error: errUnknownPopertyType("dynamic"),
		},
	}

	for input, expected := range tests {
		var (
			ctx      = logger.WithLogger(broker.NewBackground())
			property = specs.Property{}
			err      = SetDefaultValue(ctx, &property, input)
		)

		if expected.error != nil {
			if !errors.Is(err, expected.error) {
				t.Errorf("error '%s' was expected to be '%s'", err, expected.error)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error '%s'", err)
			}
		}

		if expected.defaultValue != property.Default {
			t.Errorf("unexpected result %+v, expected %+v", property.Default, expected.defaultValue)
		}

		if expected.dataType != property.Type {
			t.Errorf("unexpected type %s, expected %s", property.Type, expected.dataType)
		}
	}
}
