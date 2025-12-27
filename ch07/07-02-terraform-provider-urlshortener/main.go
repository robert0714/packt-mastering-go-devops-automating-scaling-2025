package main

import (
	"context"
	"flag"
	"log"
	customProvider "terraform-provider-urlshortener/internal"
	provider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version = "dev"
	debug   = flag.Bool("debug", false, "enable debug mode")
)

func main() {
	log.Println("Provider starting...")
	flag.Parse()
	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return customProvider.New(version)
	}, providerserver.ServeOpts{
		Address: "registry.terraform.io/local/urlshortener",
		Debug: *debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
