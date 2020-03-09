package hcl

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	log "github.com/sirupsen/logrus"
)

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
