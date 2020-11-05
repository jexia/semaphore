package graphql

import (
	"encoding/json"
	"fmt"
	"os"

	gql "github.com/graphql-go/graphql"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/config"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	print "github.com/jexia/semaphore/cmd/semaphore/generate/printer"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/endpoints"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/transport/graphql"
	"github.com/spf13/cobra"
)

var (
	flags = &config.Daemon{}

	// Command represents the semaphore daemon command
	Command = &cobra.Command{
		Use:   "graphql",
		Short: "Generates a graphql specification",
		Long: `Generates a graphql specification.
*NOTE*: This is a experimental feature and not all features are supported yet.`,
		RunE:         run,
		SilenceUsage: true,
	}

	options = print.Options{
		LineStart: "# ",
		LineEnd:   "\n",
	}
)

func init() {
	Command.PersistentFlags().StringSliceVar(&flags.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&flags.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&flags.LogLevel, "level", "warn", "Global logging level, this value will override the defined log level inside the file definitions")
}

func run(cmd *cobra.Command, args []string) (err error) {

	defer func() {
		if err != nil {
			err = prettyerr.StandardErr(err)
		}
	}()

	ctx := logger.WithLogger(broker.NewContext())
	err = config.SetOptions(ctx, flags)
	if err != nil {
		return err
	}

	core, err := config.NewCore(ctx, flags)
	if err != nil {
		return err
	}

	provider, err := config.NewProviders(ctx, core, flags)
	if err != nil {
		return err
	}

	collection, err := providers.Resolve(ctx, functions.Collection{}, provider)
	if err != nil {
		return err
	}

	transporters, err := endpoints.Transporters(ctx, collection.EndpointList, collection.FlowListInterface,
		endpoints.WithCore(core),
		endpoints.WithServices(collection.ServiceList),
		endpoints.WithFunctions(functions.Collection{}),
	)
	if err != nil {
		return err
	}

	listener := graphql.Listener{}
	err = listener.Handle(ctx, transporters, nil)
	if err != nil {
		return err
	}

	if err := print.DefaultHeader(cmd.Version).Print(os.Stdout, options); err != nil {
		return fmt.Errorf("failed to write file header: %w", err)
	}

	var (
		schema   = listener.Schema()
		params   = gql.Params{Schema: schema, RequestString: query}
		response = gql.Do(params)
		encoder  = json.NewEncoder(os.Stdout)
	)

	encoder.SetIndent("", "  ")

	return encoder.Encode(response)
}
