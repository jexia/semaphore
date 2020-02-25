package main

import (
	"fmt"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/schema/protoc"
)

func main() {
	collection, err := protoc.Collect(nil, ".")
	if err != nil {
		panic(err)
	}

	client, err := maestro.New(maestro.WithPath(".", false), maestro.WithSchemaCollection(collection))
	if err != nil {
		panic(err)
	}

	for _, flow := range client.Manifest.Flows {
		fmt.Println("flow:", flow.Name)
		for _, call := range flow.Calls {
			fmt.Println("  - call", call.Name)

			for _, prop := range call.GetRequest().GetProperties() {
				fmt.Println("    - ", prop.GetPath(), prop.Reference)
			}
		}
	}
}
