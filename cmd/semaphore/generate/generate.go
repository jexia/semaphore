package generate

import (
	"github.com/jexia/semaphore/cmd/semaphore/config"
	"github.com/jexia/semaphore/cmd/semaphore/generate/openapi3"
	"github.com/spf13/cobra"
)

var params = config.New()

// Command represents the semaphore daemon command
var Command = &cobra.Command{
	Use:          "generate",
	Short:        "Generates schemas out of the given flow and schema configurations",
	SilenceUsage: true,
}

func init() {
	Command.AddCommand(openapi3.Command)
}
