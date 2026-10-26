package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// NewProtoV5Server returns the function that creates the server Terraform talks
// to. It serves the SDKv2 provider and the framework provider as one, so that
// resources can be moved from the first to the second one at a time. A resource
// or data source is served by whichever provider declares it.
func NewProtoV5Server(ctx context.Context, version, commit string) (func() tfprotov5.ProviderServer, error) {
	sdk := newProviderFactory(version, commit)()

	mux, err := tf5muxserver.NewMuxServer(ctx,
		sdk.GRPCProvider,
		providerserver.NewProtocol5(newFrameworkProvider(version, commit)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to combine the providers: %w", err)
	}

	return mux.ProviderServer, nil
}
