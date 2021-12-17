package main

import (
	"bytes"
	"flag"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func main() {
	var provider string
	flag.StringVar(&provider, "provider", "", "provider to remove")

	flag.Parse()

	var dir string
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	} else {
		dir = "."
	}

	if provider == "" {
		panic("provider name required")
	}

	files, diags := dirFiles(tfconfig.NewOsFs(), dir)
	if diags.HasErrors() {
		panic(diags.Error())
	}

	for _, path := range files {
		original, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		file, err := hclwrite.ParseConfig(original, path, hcl.Pos{})
		RemoveProviderConfigs(provider, file)
		RemoveProviderRequirement(provider, file)

		updated := file.Bytes()
		if bytes.Compare(original, updated) == 0 {
			continue
		}

		err = os.WriteFile(path, updated, 0600)
		if err != nil {
			panic(err)
		}

	}

}
