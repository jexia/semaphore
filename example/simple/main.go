package main

import (
	"fmt"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
)

func main() {
	collection, err := protoc.Collect([]string{"../../", "."}, ".")
	if err != nil {
		panic(err)
	}

	listener, err := http.NewListener(":8080", specs.Options{})
	if err != nil {
		panic(err)
	}

	_, err = maestro.New(
		maestro.WithPath(".", false),
		maestro.WithSchemaCollection(collection),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
		maestro.WithListener(listener),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("listening on :8080")
	err = listener.Serve()
	if err != nil {
		panic(err)
	}
}
