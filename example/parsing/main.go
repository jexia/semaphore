package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
)

func main() {
	collection, err := protoc.Collect(nil, ".")
	if err != nil {
		panic(err)
	}

	listener, err := http.NewListener(":8080", specs.Options{})
	if err != nil {
		panic(err)
	}

	json := &json.Constructor{}

	_, err = maestro.New(maestro.WithPath(".", false), maestro.WithSchemaCollection(collection), maestro.WithCodec(json), maestro.WithCaller(&http.Caller{}), maestro.WithListener(listener))
	if err != nil {
		panic(err)
	}

	err = listener.Serve()
	if err != nil {
		panic(err)
	}
}
