package main

import (
	"os"

	"github.com/jexia/maestro/cmd/cli/validate"
	"github.com/jexia/maestro/cmd/daemon/run"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "maestro",
	Version: "0.1.0",
	Short:   "A straightforward micro-service conductor",
	Long: `A straightforward micro-service conductor.
Maestro is a tool to orchestrate requests inside your microservice architecture.
Requests could be manipulated, passed and branched to different services to be returned as a single output.`,
}

func init() {
	cmd.AddCommand(run.Cmd)
	cmd.AddCommand(validate.Cmd)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
