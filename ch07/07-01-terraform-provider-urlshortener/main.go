package main

import (
	"context"
	"flag"
	"log"
	provider "terraform-provider-urlshortener/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version = "dev"
	debug   = flag.Bool("debug", false, "enable debug mode")
)

func main() {
	log.Println("Provider starting...")
	flag.Parse()
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/local/urlshortener",
		Debug:   *debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
