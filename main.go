package main

import (
	"context"
	"flag"
	"log"

	"github.com/datadrivers/terraform-provider-nexus/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
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

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/datadrivers/nexus",
			&plugin.ServeOpts{
				ProviderFunc: provider.Provider,
			})
		if err != nil {
			log.Printf("[ERROR] Error during initialization: %s", err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: provider.Provider})
	}
}
