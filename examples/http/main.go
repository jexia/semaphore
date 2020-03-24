package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport/http"
)

func main() {
	collection, err := protoc.Collect([]string{"../../annotations", "."}, "./*.proto")
	if err != nil {
		panic(err)
	}

	client, err := maestro.New(
		maestro.WithLogLevel(logger.Global, "debug"),
		maestro.WithListener(http.NewListener(":8080", specs.Options{})),
		maestro.WithDefinitions(hcl.DefinitionResolver("./*.hcl")),
		maestro.WithSchema(collection),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
	)

	if err != nil {
		panic(err)
	}

	err = client.Serve()
	if err != nil {
		panic(err)
	}
}
