package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/kyma-incubator/terraform-provider-gardener/gardener"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gardener.Provider})
}
