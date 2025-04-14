// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"landb/internal/provider"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/barnes-c/terraform-provider-landb",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
