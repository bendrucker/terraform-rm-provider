package main

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
)

func TestRemoveProviderConfigs(t *testing.T) {
	cases := []struct {
		name     string
		config   string
		remove   string
		expected string
	}{
		{
			name: "basic",
			config: `
				provider "foo" {}
				provider "bar" {}
		`,
			remove: "foo",
			expected: `
				provider "bar" {}
		`,
		},
		{
			name: "noop",
			config: `
				provider "foo" {}
				provider "bar" {}
		`,
			remove: "nomatch",
			expected: `
				provider "foo" {}
				provider "bar" {}
		`,
		},
		{
			name: "aliases",
			config: `
				provider "foo" {}
				provider "foo" {
					alias = "a"
				}
				provider "bar" {}
		`,
			remove: "foo",
			expected: `
				provider "bar" {}
		`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f, err := hclwrite.ParseConfig([]byte(tc.config), "main.tf", hcl.Pos{})
			if err != nil {
				t.Fatal(err)
			}

			RemoveProviderConfigs(tc.remove, f)

			expected := strings.TrimSpace(dedent.Dedent(tc.expected))
			actual := strings.TrimSpace(string(f.Bytes()))

			assert.Equal(t, expected, actual)
		})
	}
}

func TestRemoveProviderRequirement(t *testing.T) {
	cases := []struct {
		name     string
		config   string
		remove   string
		expected string
	}{
		{
			name: "basic",
			config: `
				terraform {
					required_version = ">= 1"
					
					required_providers {
						foo = {
							source = "ns/foo"
						}

						bar = {
							source = "ns/bar"
						}
					}
				}
		`,
			remove: "foo",
			expected: `
			terraform {
				required_version = ">= 1"
				
				required_providers {

					bar = {
						source = "ns/bar"
					}
				}
			}
		`,
		},
		{
			name: "noop",
			config: `
				terraform {
					required_version = ">= 1"
					
					required_providers {
						foo = {
							source = "ns/foo"
						}

						bar = {
							source = "ns/bar"
						}
					}
				}
		`,
			remove: "nomatch",
			expected: `
				terraform {
					required_version = ">= 1"
					
					required_providers {
						foo = {
							source = "ns/foo"
						}

						bar = {
							source = "ns/bar"
						}
					}
				}
		`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f, err := hclwrite.ParseConfig([]byte(tc.config), "main.tf", hcl.Pos{})
			if err != nil {
				t.Fatal(err)
			}

			RemoveProviderRequirement(tc.remove, f)

			expected := hclwrite.Format([]byte(tc.expected))
			actual := hclwrite.Format(f.Bytes())

			assert.Equal(t, string(expected), string(actual))
		})
	}
}
