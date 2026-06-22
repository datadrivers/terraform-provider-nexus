package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/williamt1997/terraform-provider-nexus/internal/provider"
)

// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Clean up log output
	// See https://developer.hashicorp.com/terraform/plugin/log/writing#legacy-logging
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		Debug:        debugMode,
		ProviderAddr: "registry.terraform.io/williamt1997/nexus",
		ProviderFunc: provider.Provider,
	})
}
