package config

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// New constructs a new global config
func New() *Maestro {
	return &Maestro{
		HTTP:         HTTP{},
		GraphQL:      GraphQL{},
		GRPC:         GRPC{},
		Protobuffers: []string{},
		HCL:          []string{},
	}
}

// Read attempts to read the given configuration on the given path and decode it into the Maestro configuration
func Read(cmd *cobra.Command, target *Maestro) error {
	flag := cmd.Flag("config")
	if flag == nil {
		return nil
	}

	path := flag.Value.String()
	if path == "" {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	err = yaml.NewDecoder(file).Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

// Maestro configurations
type Maestro struct {
	LogLevel     string
	HTTP         HTTP
	GraphQL      GraphQL
	GRPC         GRPC
	Protobuffers []string
	HCL          []string
}

// HTTP configurations
type HTTP struct {
	Address string
}

// GRPC configurations
type GRPC struct {
	Address string
}

// GraphQL configurations
type GraphQL struct {
	Address string
}
