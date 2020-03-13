package hcl

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/utils"
	log "github.com/sirupsen/logrus"
)

// SchemaResolver constructs a schema resolver for the given path
func SchemaResolver(path string) schema.Resolver {
	recursive := false
	if strings.HasSuffix(path, "...") {
		recursive = true
		path = path[:3]
	}

	return func(schemas *schema.Store) error {
		files, err := utils.ReadDir(path, recursive, Ext)
		if err != nil {
			return err
		}

		for _, file := range files {
			reader, err := os.Open(filepath.Join(file.Path, file.Name()))
			if err != nil {
				return err
			}

			definition, err := UnmarshalHCL(file.Name(), reader)
			if err != nil {
				return err
			}

			collection, err := ParseSchema(definition, schemas)
			if err != nil {
				return err
			}

			schemas.Add(collection)
		}

		return nil
	}
}

// DefinitionResolver constructs a definition resolver for the given path
func DefinitionResolver(path string) specs.Resolver {
	recursive := false
	if strings.HasSuffix(path, "...") {
		recursive = true
		path = path[:3]
	}

	return func(functions specs.CustomDefinedFunctions) (*specs.Manifest, error) {
		files, err := utils.ReadDir(path, recursive, Ext)
		if err != nil {
			return nil, err
		}

		result := &specs.Manifest{}

		for _, file := range files {
			reader, err := os.Open(filepath.Join(file.Path, file.Name()))
			if err != nil {
				return nil, err
			}

			definition, err := UnmarshalHCL(file.Name(), reader)
			if err != nil {
				return nil, err
			}

			manifest, err := ParseSpecs(definition, functions)
			if err != nil {
				return nil, err
			}

			result.Merge(manifest)
		}

		return result, nil
	}
}

// UnmarshalHCL unmarshals the given HCL stream into a intermediate resource.
func UnmarshalHCL(filename string, reader io.Reader) (manifest Manifest, _ error) {
	log.WithField("file", filename).Info("Reading HCL files")

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return manifest, err
	}

	log.WithField("file", filename).Debug("Parsing HCL syntax")

	file, diags := hclsyntax.ParseConfig(bb, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	log.WithField("file", filename).Debug("Decoding HCL syntax")

	diags = gohcl.DecodeBody(file.Body, nil, &manifest)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	return manifest, nil
}
