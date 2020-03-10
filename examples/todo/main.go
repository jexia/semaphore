package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/sirupsen/logrus"
)

func main() {
	collection, err := protoc.Collect([]string{"../../", "."}, ".")
	if err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.DebugLevel)

	listener, err := http.NewListener(":8080", specs.Options{})
	if err != nil {
		panic(err)
	}

	_, err = maestro.New(
		maestro.WithPath(".", false),
		maestro.WithSchema(collection),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
		maestro.WithListener(listener),
	)

	if err != nil {
		panic(err)
	}

	err = listener.Serve()
	if err != nil {
		panic(err)
	}
}
