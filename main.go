package main

import (
	"github.com/datadrivers/terraform-provider-nexus/nexus"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nexus.Provider(),
	})
}
