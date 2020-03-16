package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol/graphql"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/sirupsen/logrus"
)

func main() {
	collection, err := protoc.Collect([]string{"../../annotations", "."}, "./*.proto")
	if err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.DebugLevel)

	client, err := maestro.New(
		maestro.WithDefinitions(hcl.DefinitionResolver("./*.hcl")),
		maestro.WithSchema(collection),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
		maestro.WithListener(graphql.NewListener(":8080", specs.Options{})),
	)

	if err != nil {
		panic(err)
	}

	err = client.Serve()
	if err != nil {
		panic(err)
	}
}
