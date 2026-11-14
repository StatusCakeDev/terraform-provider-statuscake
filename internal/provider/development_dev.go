//go:build dev

package provider

import (
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// descCustomEndpoint describes the setting for both providers.
const descCustomEndpoint = "Custom endpoint to which requests will be made. Only available in development builds. This can also be provided as an environment variable `STATUSCAKE_CUSTOM_ENDPOINT`"

// developmentSchema returns the settings that exist only in development builds,
// which are made with `go build -tags dev` (see `make build TAGS=dev`).
func developmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_endpoint": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_CUSTOM_ENDPOINT", nil),
			Description:  descCustomEndpoint,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

// developmentFrameworkAttributes is developmentSchema for the framework
// provider.
func developmentFrameworkAttributes() map[string]fwschema.Attribute {
	return map[string]fwschema.Attribute{
		"custom_endpoint": fwschema.StringAttribute{
			Optional:    true,
			Description: descCustomEndpoint,
		},
	}
}

// developmentEndpoint returns the custom endpoint that was configured, or an
// empty string when there is none.
func developmentEndpoint(d *schema.ResourceData) string {
	endpoint, _ := d.Get("custom_endpoint").(string)
	return endpoint
}
