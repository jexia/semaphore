package main

import (
	"errors"
	"fmt"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
	"github.com/jexia/maestro/transport/http"
)

func main() {
	collection, err := protoc.Collect([]string{"../../annotations", "./proto"}, "./proto/*.proto")
	if err != nil {
		panic(err)
	}

	functions := specs.CustomDefinedFunctions{
		"jwt": jwt,
	}

	client, err := maestro.New(
		maestro.WithLogLevel(logger.Global, "debug"),
		maestro.WithListener(http.NewListener(":8080", specs.Options{})),
		maestro.WithDefinitions(hcl.DefinitionResolver("./*.hcl")),
		maestro.WithSchema(collection),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
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

func jwt(args ...*specs.Property) (*specs.Property, specs.FunctionExec, error) {
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

	fn := func(store specs.Store) error {
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
