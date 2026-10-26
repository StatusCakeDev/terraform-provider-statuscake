// Package provider implements the Terraform provider for StatusCake.
package provider

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/adapters"
)

// Descriptions of the provider's settings. Both the SDKv2 provider and the
// framework provider that is served beside it declare these settings, and
// Terraform requires their schemas to be identical.
const (
	descAPIToken   = "The API token for operations. This can also be provided as an environment variable `STATUSCAKE_API_TOKEN`"
	descRPS        = "RPS limit to apply when making calls to the API. This can also be provided as an environment variable `STATUSCAKE_RPS`"
	descRetries    = "Maximum number of retries to perform when an API request fails. This can also be provided as an environment variable `STATUSCAKE_RETRIES`"
	descMinBackoff = "Minimum backoff period in seconds after failed API calls. This can also be provided as an environment variable `STATUSCAKE_MIN_BACKOFF`"
	descMaxBackoff = "Maximum backoff period in seconds after failed API calls. This can also be provided as an environment variable `STATUSCAKE_MAX_BACKOFF`"
)

// newProviderFactory returns the function that creates the SDKv2 provider. The
// version and commit are used to set the version of the provider in the user
// agent string.
func newProviderFactory(version, commit string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_token": {
					Type:         schema.TypeString,
					Required:     true,
					DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_API_TOKEN", nil),
					Description:  descAPIToken,
					ValidateFunc: validation.StringMatch(regexp.MustCompile("[0-9a-zA-Z_]{20,30}"), "API token must only contain characters 0-9, a-zA-Z and underscores"),
				},
				"rps": {
					Type:         schema.TypeInt,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_RPS", 4),
					Description:  descRPS,
					ValidateFunc: validation.IntAtLeast(1),
				},
				"retries": {
					Type:         schema.TypeInt,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_RETRIES", 3),
					Description:  descRetries,
					ValidateFunc: validation.IntBetween(0, 10),
				},
				"min_backoff": {
					Type:         schema.TypeInt,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_MIN_BACKOFF", 1),
					Description:  descMinBackoff,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"max_backoff": {
					Type:         schema.TypeInt,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_MAX_BACKOFF", 30),
					Description:  descMaxBackoff,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"statuscake_contact_group":      resourceStatusCakeContactGroup(),
				"statuscake_heartbeat_check":    resourceStatusCakeHeartbeatCheck(),
				"statuscake_maintenance_window": resourceStatusCakeMaintenanceWindow(),
				"statuscake_pagespeed_check":    resourceStatusCakePagespeedCheck(),
				"statuscake_ssl_check":          resourceStatusCakeSSLCheck(),
				"statuscake_uptime_check":       resourceStatusCakeUptimeCheck(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"statuscake_contact_group":                  dataSourceStatusCakeContactGroup(),
				"statuscake_pagespeed_monitoring_locations": dataSourceStatusCakeMonitoringLocations(listPagespeedMonitoringLocations),
				"statuscake_uptime_monitoring_locations":    dataSourceStatusCakeMonitoringLocations(listUptimeMonitoringLocations),
			},
			ConfigureContextFunc: providerConfigure(version, commit),
		}

		// Settings that exist only in development builds.
		for name, s := range developmentSchema() {
			p.Schema[name] = s
		}

		return p
	}
}

// providerConfigure parses the config into the Terraform provider meta object.
func providerConfigure(version, commit string) func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiToken, ok := d.GetOk("api_token")
		if !ok {
			return nil, diag.Errorf("credentials are not set correctly")
		}

		client, err := adapters.New(adapters.Config{
			Token:      apiToken.(string),
			RPS:        d.Get("rps").(int),
			Retries:    d.Get("retries").(int),
			MinBackoff: time.Duration(d.Get("min_backoff").(int)) * time.Second,
			MaxBackoff: time.Duration(d.Get("max_backoff").(int)) * time.Second,
			Endpoint:   developmentEndpoint(d),
			UserAgent:  userAgent(version, commit),
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}

func userAgent(version, commit string) string {
	if version == "" {
		version = "dev"
	}
	if commit != "" {
		version = fmt.Sprintf("%s (%s)", version, commit)
	}
	return fmt.Sprintf("terraform-provider-statuscake/%s (%s; %s; %s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
