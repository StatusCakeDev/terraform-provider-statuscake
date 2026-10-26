package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// frameworkProvider is the provider as built with the Plugin Framework. It is
// served beside the SDKv2 provider (see NewProtoV5Server), and owns the
// resources and data sources that have been ported to the framework. Every
// other one is still served by the SDKv2 provider.
//
// Terraform configures the provider once, but both servers receive the
// configuration. For this reason the settings are declared by both, and must be
// identical. A test asserts that they are.
type frameworkProvider struct {
	version string
	commit  string
}

var _ fwprovider.Provider = (*frameworkProvider)(nil)

func newFrameworkProvider(version, commit string) *frameworkProvider {
	return &frameworkProvider{version: version, commit: commit}
}

func (p *frameworkProvider) Metadata(_ context.Context, _ fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "statuscake"
	resp.Version = p.version
}

func (p *frameworkProvider) Schema(_ context.Context, _ fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	// The SDKv2 provider makes api_token required unless the environment variable
	// that can supply it is set, as the schema it sends to Terraform treats an
	// attribute with a default as optional.
	token, ok := os.LookupEnv("STATUSCAKE_API_TOKEN")
	tokenInEnvironment := ok && token != ""

	attributes := map[string]fwschema.Attribute{
		"api_token": fwschema.StringAttribute{
			Required:    !tokenInEnvironment,
			Optional:    tokenInEnvironment,
			Description: descAPIToken,
		},
		"rps": fwschema.Int64Attribute{
			Optional:    true,
			Description: descRPS,
		},
		"retries": fwschema.Int64Attribute{
			Optional:    true,
			Description: descRetries,
		},
		"min_backoff": fwschema.Int64Attribute{
			Optional:    true,
			Description: descMinBackoff,
		},
		"max_backoff": fwschema.Int64Attribute{
			Optional:    true,
			Description: descMaxBackoff,
		},
	}

	// Settings that exist only in development builds.
	for name, a := range developmentFrameworkAttributes() {
		attributes[name] = a
	}

	resp.Schema = fwschema.Schema{Attributes: attributes}
}

// Configure does nothing yet, as there is nothing the framework serves that
// needs the API client. The SDKv2 provider validates the settings and builds
// the client for the resources it serves.
//
// When the first resource or data source is ported this reads the settings,
// including their defaults from the environment, and builds the same client for
// them. The settings then also need to be validated here, which is what the
// SDKv2 provider does for them now. Doing so earlier would report each invalid
// setting twice.
func (p *frameworkProvider) Configure(context.Context, fwprovider.ConfigureRequest, *fwprovider.ConfigureResponse) {
}

func (p *frameworkProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *frameworkProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}
