package main

import (
	"fmt"
	"os"

	"github.com/jexia/maestro/cmd/maestro/daemon"
	"github.com/jexia/maestro/cmd/maestro/validate"
	"github.com/spf13/cobra"
)

var version string
var build string

var cmd = &cobra.Command{
	Use:     "maestro",
	Version: fmt.Sprintf("%s, build: %s", version, build),
	Short:   "A straightforward micro-service conductor",
	Long: `Maestro is a feature-rich service orchestrator.
Create advanced data flows and expose them through endpoints. Have full control over your exposed endpoints,
expose single flows for multiple protocols such as gRPC and GraphQL. Maestro adapts to your environment.
Create custom extensions or use the availability of custom functions and protocol implementations.`,
}

func init() {
	cmd.AddCommand(daemon.Command)
	cmd.AddCommand(validate.Command)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
