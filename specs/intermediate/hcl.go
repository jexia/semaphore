package intermediate

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// UnmarshalHCL unmarshals the given HCL stream into a intermediate resource.
func UnmarshalHCL(filename string, reader io.Reader) (manifest Manifest, _ error) {
	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return manifest, err
	}

	file, diags := hclsyntax.ParseConfig(bb, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	diags = gohcl.DecodeBody(file.Body, nil, &manifest)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	return manifest, nil
}
