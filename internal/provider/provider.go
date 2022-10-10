package provider

import (
	"context"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"golang.org/x/time/rate"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/StatusCakeDev/statuscake-go/backoff"
	"github.com/StatusCakeDev/statuscake-go/credentials"
	"github.com/StatusCakeDev/statuscake-go/throttle"
)

// Provider returns a resource provider for Terraform.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_API_TOKEN", nil),
				Description:  "The API token for operations. This can also be provided as an environment variable `STATUSCAKE_API_TOKEN`",
				ValidateFunc: validation.StringMatch(regexp.MustCompile("[0-9a-zA-Z_]{20,30}"), "API token must only contain characters 0-9, a-zA-Z and underscores"),
			},
			"rps": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_RPS", 4),
				Description:  "RPS limit to apply when making calls to the API. This can also be provided as an environment variable `STATUSCAKE_RPS`",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"retries": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_RETRIES", 3),
				Description:  "Maximum number of retries to perform when an API request fails. This can also be provided as an environment variable `STATUSCAKE_RETRIES`",
				ValidateFunc: validation.IntBetween(0, 10),
			},
			"min_backoff": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_MIN_BACKOFF", 1),
				Description:  "Minimum backoff period in seconds after failed API calls. This can also be provided as an environment variable `STATUSCAKE_MIN_BACKOFF`",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_backoff": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_MAX_BACKOFF", 30),
				Description:  "Maximum backoff period in seconds after failed API calls. This can also be provided as an environment variable `STATUSCAKE_MAX_BACKOFF`",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"statuscake_custom_endpoint": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUCAKE_CUSTOM_ENDPOINT", nil),
				Description:  "Custom endpoint to which request will be made. This can also be provided as an environment variable `STATUCAKE_CUSTOM_ENDPOINT`",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"statuscake_contact_group":      resourceStatusCakeContactGroup(),
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
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure parses the config into the Terraform provider meta object.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiToken, ok := d.GetOk("api_token")
	if !ok {
		return nil, diag.Errorf("credentials are not set correctly")
	}

	bearer := credentials.NewBearerWithStaticToken(apiToken.(string))
	opts := []statuscake.Option{
		statuscake.WithBackoff(backoff.Exponential{
			BaseDelay:  time.Duration(d.Get("min_backoff").(int)) * time.Second,
			Multiplier: 2.0,
			Jitter:     0.2,
			MaxDelay:   time.Duration(d.Get("max_backoff").(int)) * time.Second,
		}),
		statuscake.WithHTTPClient(&http.Client{
			Transport: throttle.NewWithDefaultTransport(
				rate.NewLimiter(rate.Limit(d.Get("rps").(int)), 1),
			),
		}),
		statuscake.WithMaxRetries(d.Get("retries").(int)),
		statuscake.WithRequestCredentials(bearer),
		statuscake.WithUserAgent("terraform-provider-statuscake/" + runtime.Version()),
	}

	if customEndpoint, ok := d.GetOk("statuscake_custom_endpoint"); ok {
		opts = append(opts, statuscake.WithHost(customEndpoint.(string)))
	}

	return statuscake.NewClient(opts...), nil
}
