package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	demo "terraform-provider-yunjidemo/demo"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: demo.Provider,
	})
}
