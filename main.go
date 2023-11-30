package main

import (
	"context"
	"terraform-provider-misc/misc"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name misc

func main() {
	providerserver.Serve(context.Background(), misc.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/kiwicom/kiwi",
	})
}
