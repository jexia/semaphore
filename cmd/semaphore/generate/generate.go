package generate

import (
	"github.com/jexia/semaphore/v2/cmd/semaphore/generate/graphql"
	"github.com/jexia/semaphore/v2/cmd/semaphore/generate/openapi3"
	"github.com/jexia/semaphore/v2/cmd/semaphore/generate/protobuf"

	"github.com/spf13/cobra"
)

// Command represents the semaphore daemon command
var Command = &cobra.Command{
	Use:          "generate",
	Short:        "Generates schemas out of the given flow and schema configurations",
	SilenceUsage: true,
}

func init() {
	Command.AddCommand(openapi3.Command)
	Command.AddCommand(protobuf.Command)
	Command.AddCommand(graphql.Command)
}
