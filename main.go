package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-warpgate/provider"

	// "terraform-provider-warpgate/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "0.1.0"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		Address: "registry.terraform.io/andreee94/warpgate",
		// Address: "local/tr/warpgate",
		Debug: debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}