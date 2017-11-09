package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ustream/terraform-provider-compose/compose"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: compose.Provider,
	})
}
