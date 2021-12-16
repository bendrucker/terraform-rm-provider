package main

import "github.com/hashicorp/hcl/v2/hclwrite"

func RemoveProviderConfigs(name string, file *hclwrite.File) {
	body := file.Body()

	for _, block := range body.Blocks() {
		if block.Type() != "provider" || block.Labels()[0] != name {
			continue
		}

		body.RemoveBlock(block)
	}
}

func RemoveProviderRequirement(name string, file *hclwrite.File) {
	body := file.Body()

	for _, block := range body.Blocks() {
		if block.Type() != "terraform" {
			continue
		}

		for _, block := range block.Body().Blocks() {
			if block.Type() != "required_providers" {
				continue
			}

			for key, _ := range block.Body().Attributes() {
				if key != name {
					continue
				}

				block.Body().RemoveAttribute(key)
			}
		}
	}
}
