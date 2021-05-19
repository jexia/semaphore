package main

import (
	"fmt"
	"os"

	"github.com/jexia/semaphore/v2/cmd/semaphore/daemon"
	"github.com/jexia/semaphore/v2/cmd/semaphore/generate"
	"github.com/jexia/semaphore/v2/cmd/semaphore/validate"
	"github.com/spf13/cobra"
)

var (
	version string
	build   string
	label   string
)

var cmd = &cobra.Command{
	Use:     "semaphore",
	Version: fmt.Sprintf("%s, build: %s %s", version, build, label),
	Short:   "A straightforward micro-service conductor",
	Long: `Semaphore is a feature-rich service orchestrator.
Create advanced data flows and expose them through endpoints. Have full control over your exposed endpoints,
expose single flows for multiple protocols such as gRPC and GraphQL. Semaphore adapts to your environment.
Create custom extensions or use the availability of custom functions and protocol implementations.`,
}

func init() {
	cmd.AddCommand(daemon.Command)
	cmd.AddCommand(validate.Command)
	cmd.AddCommand(generate.Command)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
