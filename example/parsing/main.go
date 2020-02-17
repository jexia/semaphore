package main

import (
	"fmt"
	"path/filepath"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/schema/protoc"
)

func main() {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	collection, err := protoc.Collect([]string{path}, path)
	if err != nil {
		panic(err)
	}

	manifest, err := maestro.New(maestro.WithPath(".", false), maestro.WithSchemaCollection(collection))
	if err != nil {
		panic(err)
	}

	for _, flow := range manifest.Flows {
		fmt.Println("flow:", flow.Name)
		for _, call := range flow.Calls {
			fmt.Println("  - call", call.Name)
		}
	}
}
