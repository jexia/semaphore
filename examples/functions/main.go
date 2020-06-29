package main

import (
	"errors"
	"fmt"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/internal/functions"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/providers/hcl"
	"github.com/jexia/maestro/pkg/providers/proto"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jexia/maestro/pkg/transport/http"
)

func main() {
	functions := functions.Custom{
		"jwt": jwt,
	}

	client, err := maestro.New(
		maestro.WithLogLevel(logger.Global, "debug"),
		maestro.WithListener(http.NewListener(":8080", specs.Options{})),
		maestro.WithFlows(hcl.FlowsResolver("./*.hcl")),
		maestro.WithEndpoints(hcl.EndpointsResolver("./*.hcl")),
		maestro.WithSchema(proto.SchemaResolver([]string{"../../../", "./proto"}, "./proto/*.proto")),
		maestro.WithServices(proto.ServiceResolver([]string{"../../../", "./proto"}, "./proto/*.proto")),
		maestro.WithCodec(codec.JSON()),
		maestro.WithCodec(codec.Proto()),
		maestro.WithCaller(http.NewCaller()),
		maestro.WithFunctions(functions),
	)

	if err != nil {
		panic(err)
	}

	err = client.Serve()
	if err != nil {
		panic(err)
	}
}

func jwt(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
	prop := &specs.Property{
		Type:  types.String,
		Label: labels.Optional,
	}

	if len(args) != 1 {
		return nil, nil, fmt.Errorf("invalid jwt amount of arguments %d, expected 1", len(args))
	}

	input := args[0]

	if input.Type != types.String {
		return nil, nil, fmt.Errorf("invalid argument type (%s), expected (%s)", input.Type, types.String)
	}

	fn := func(store refs.Store) error {
		value := input.Default

		if input.Reference != nil {
			ref := store.Load(input.Reference.Resource, input.Reference.Path)
			if ref != nil {
				value = ref.Value
			}
		}

		if value == nil {
			return errors.New("invalid token")
		}

		token, is := value.(string)
		if !is {
			return errors.New("invalid value, expected a string")
		}

		if token != "super-secret" {
			return errors.New("token is invalid expected 'super-secret'")
		}

		return nil
	}

	return prop, fn, nil
}
