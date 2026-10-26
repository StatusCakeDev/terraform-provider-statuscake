//go:build !dev

package provider

import (
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// developmentSchema returns the settings that exist only in development builds.
// A release build has none, so users cannot point the provider at anything but
// the StatusCake API.
func developmentSchema() map[string]*schema.Schema { return nil }

// developmentFrameworkAttributes is developmentSchema for the framework provider.
func developmentFrameworkAttributes() map[string]fwschema.Attribute { return nil }

// developmentEndpoint is always empty in a release build, which has no setting
// to read.
func developmentEndpoint(*schema.ResourceData) string { return "" }
