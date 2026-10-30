package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for
// the registry/website

// If you do not have terraform installed, you can remove the formatting
// command, but its suggested to ensure the documentation is formatted properly.
//
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on
// how it works and how docs can be customized.
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// Build information of the server, set at build time via ldflags.
	version = "dev"
	commit  = ""
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// The SDKv2 provider and the framework provider are served as one. The
	// resources are moved from the first to the second one at a time.
	server, err := provider.NewProtoV5Server(context.Background(), version, commit)
	if err != nil {
		log.Fatal(err)
	}

	opts := &plugin.ServeOpts{
		GRPCProviderFunc: server,
		Debug:            debug,
		ProviderAddr:     "registry.terraform.io/StatusCakeDev/statuscake",
	}

	plugin.Serve(opts)
}
