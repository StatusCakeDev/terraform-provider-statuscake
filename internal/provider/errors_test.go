//go:build dev

package provider_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAPIViolationsAreReportedAgainstTheField(t *testing.T) {
	api := newFakeAPI(t)
	api.intercept = func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodPost {
			return false
		}
		writeError(w, http.StatusBadRequest, "Bad request", map[string][]string{
			"name": {"is already in use"},
		})
		return true
	}

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_contact_group" "test" {
					name = "Operations"
				}
			`,
			ExpectError: regexp.MustCompile("(?s)failed to create contact group.*name contains violations.*is already in use"),
		},
	)
}

func TestAPIFailureIsReportedWhenReading(t *testing.T) {
	api := newFakeAPI(t)
	id := api.Seed("contact-groups", map[string][]string{"name": {"Operations"}})
	api.intercept = func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodGet {
			return false
		}
		writeError(w, http.StatusInternalServerError, "Internal error", nil)
		return true
	}

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				data "statuscake_contact_group" "test" {
					id = "` + id + `"
				}
			`,
			ExpectError: regexp.MustCompile("failed to get contact group with ID: " + id + ", error"),
		},
	)
}

// A resource that has been deleted outside of Terraform is dropped from the
// state when it is refreshed, and so is created again by the next apply.
func TestResourceDeletedOutsideTerraformIsRecreated(t *testing.T) {
	tests := map[string]struct {
		kind   string
		config string
	}{
		"statuscake_contact_group": {
			kind:   "contact-groups",
			config: `name = "Operations"`,
		},
		"statuscake_heartbeat_check": {
			kind: "heartbeat",
			config: `
				name   = "Job"
				period = 1800
			`,
		},
		"statuscake_maintenance_window": {
			kind: "maintenance-windows",
			config: `
				name     = "Window"
				start    = "2022-01-29T00:00:00Z"
				end      = "2022-01-30T23:59:59Z"
				timezone = "UTC"
				tags     = ["production"]
			`,
		},
		"statuscake_pagespeed_check": {
			kind: "pagespeed",
			config: `
				check_interval = 300
				name           = "Example"
				region         = "UK"

				alert_config {}

				monitored_resource {
					address = "https://www.example.com"
				}
			`,
		},
		"statuscake_ssl_check": {
			kind: "ssl",
			config: `
				check_interval = 600

				alert_config {
					alert_at = [1, 7, 30]
				}

				monitored_resource {
					address = "https://www.example.com"
				}
			`,
		},
		"statuscake_uptime_check": {
			kind: "uptime",
			config: `
				check_interval = 300
				name           = "Example"

				http_check {}

				monitored_resource {
					address = "https://www.example.com"
				}
			`,
		},
	}
	for resourceType, tt := range tests {
		t.Run(resourceType, func(t *testing.T) {
			api := newFakeAPI(t)
			config := providerBlock() + `
				resource "` + resourceType + `" "test" {
					` + tt.config + `
				}
			`

			runTest(t, api,
				resource.TestStep{Config: config},
				resource.TestStep{
					PreConfig: func() { api.Clear(tt.kind) },
					Config:    config,
				},
			)

			if got := len(api.Calls("POST", "/v1/"+tt.kind)); got != 2 {
				t.Errorf("got %d creates, want 2: the resource should have been created again", got)
			}
		})
	}
}

func TestProviderRejectsAMalformedToken(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: `
				provider "statuscake" {
					api_token = "short"
				}

				data "statuscake_uptime_monitoring_locations" "test" {}
			`,
			ExpectError: regexp.MustCompile("API token must only contain characters"),
		},
	)
}

func TestProviderSendsItsUserAgentAndToken(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `data "statuscake_uptime_monitoring_locations" "test" {}`,
		},
	)

	calls := api.Calls("GET", "/v1/uptime-locations")
	if len(calls) == 0 {
		t.Fatal("no request was made")
	}

	if got, want := calls[0].Header.Get("Authorization"), "Bearer "+testToken; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	if got := calls[0].Header.Get("User-Agent"); !regexp.MustCompile(`^terraform-provider-statuscake/test \(abc123\) \(go[^;]*; \w+; \w+\)$`).MatchString(got) {
		t.Errorf("User-Agent = %q", got)
	}
}
