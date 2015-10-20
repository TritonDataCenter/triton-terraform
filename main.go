package main

import (
	"github.com/hashicorp/terraform/plugin"
)

// FakePlugin is a fake plugin
type FakePlugin struct{}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
