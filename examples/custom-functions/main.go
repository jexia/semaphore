package main

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport/http"
)

func main() {
	functions := functions.Custom{
		"jwt": jwt,
	}

	client, err := semaphore.New(
		semaphore.WithLogLevel("*", "debug"),
		semaphore.WithListener(http.NewListener(":8080", specs.Options{})),
		semaphore.WithFlows(hcl.FlowsResolver("./*.hcl")),
		semaphore.WithEndpoints(hcl.EndpointsResolver("./*.hcl")),
		semaphore.WithSchema(protobuffers.SchemaResolver([]string{"./proto"}, "./proto/*.proto")),
		semaphore.WithServices(protobuffers.ServiceResolver([]string{"./proto"}, "./proto/*.proto")),
		semaphore.WithCodec(json.NewConstructor()),
		semaphore.WithCodec(proto.NewConstructor()),
		semaphore.WithCaller(http.NewCaller()),
		semaphore.WithFunctions(functions),
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

	fn := func(store references.Store) error {
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
