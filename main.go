package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/StatusCakeDev/terraform-provider-statuscake/internal/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.Provider}

	if debug {
		if err := plugin.Debug(context.Background(), "registry.terraform.io/StatusCakeDev/statuscake", opts); err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)
}
