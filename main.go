package main

import (
	"github.com/hashicorp/terraform/plugin"
)

const Name = "triton-terraform"
const Version = "0.0.1"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
