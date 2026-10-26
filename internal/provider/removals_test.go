//go:build dev

package provider_test

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Removing an optional setting must blank it on the API, because the API keeps
// any value that it is not sent. Each case applies a configuration with the
// settings, then one without them, and checks that the update sent a blank
// value for every removed setting.
func TestRemovingOptionalSettingsBlanksThemOnTheAPI(t *testing.T) {
	tests := map[string]struct {
		// full and bare are the configuration of the resource's body.
		full string
		bare string
		// cleared maps each field expected in the update to the values sent.
		cleared url.Values
	}{
		"statuscake_contact_group": {
			full: `
				name            = "Operations"
				email_addresses = ["ops@example.com"]
				mobile_numbers  = ["+441234567890"]
				integrations    = ["1"]
				ping_url        = "https://example.com/hook"
			`,
			bare: `name = "Operations"`,
			cleared: url.Values{
				"email_addresses[]": {""},
				"mobile_numbers[]":  {""},
				"integrations[]":    {""},
				"ping_url":          {""},
			},
		},
		"statuscake_heartbeat_check": {
			full: `
				name           = "Example"
				period         = 60
				tags           = ["a"]
				contact_groups = ["1"]
				monitored_resource {
					host = "Example Hosting"
				}
			`,
			bare: `
				name   = "Example"
				period = 60
			`,
			cleared: url.Values{
				"tags[]":           {""},
				"contact_groups[]": {""},
				"host":             {""},
			},
		},
		"statuscake_maintenance_window": {
			full: `
				name     = "Example"
				start    = "2030-01-01T00:00:00Z"
				end      = "2030-01-01T01:00:00Z"
				timezone = "UTC"
				tags     = ["a"]
				tests    = ["1"]
			`,
			bare: `
				name     = "Example"
				start    = "2030-01-01T00:00:00Z"
				end      = "2030-01-01T01:00:00Z"
				timezone = "UTC"
				tags     = ["a"]
			`,
			cleared: url.Values{"tests[]": {""}},
		},
		"statuscake_ssl_check": {
			full: `
				check_interval = 600
				user_agent     = "agent"
				contact_groups = ["1"]
				alert_config {
					alert_at = [1, 7, 30]
				}
				monitored_resource {
					address  = "https://www.example.com"
					hostname = "www.example.com"
				}
			`,
			bare: `
				check_interval = 600
				alert_config {
					alert_at = [1, 7, 30]
				}
				monitored_resource {
					address = "https://www.example.com"
				}
			`,
			cleared: url.Values{
				"user_agent":       {""},
				"contact_groups[]": {""},
				"hostname":         {""},
			},
		},
		"statuscake_pagespeed_check": {
			full: `
				check_interval = 300
				name           = "Example"
				region         = "UK"
				contact_groups = ["1"]
				alert_config {}
				monitored_resource {
					address = "https://www.example.com"
				}
			`,
			bare: `
				check_interval = 300
				name           = "Example"
				region         = "UK"
				alert_config {}
				monitored_resource {
					address = "https://www.example.com"
				}
			`,
			cleared: url.Values{"contact_groups[]": {""}},
		},
		"statuscake_uptime_check": {
			full: `
				check_interval = 300
				name           = "Example"
				tags           = ["a"]
				contact_groups = ["1"]
				regions        = ["london"]
				http_check {
					final_endpoint      = "https://www.example.com/end"
					user_agent          = "agent"
					request_payload_raw = "payload"
					request_headers     = { a = "b" }
				}
				monitored_resource {
					address = "https://www.example.com"
					host    = "Example Hosting"
				}
			`,
			bare: `
				check_interval = 300
				name           = "Example"
				http_check {}
				monitored_resource {
					address = "https://www.example.com"
				}
			`,
			cleared: url.Values{
				"tags[]":           {""},
				"contact_groups[]": {""},
				"regions[]":        {""},
				"final_endpoint":   {""},
				"user_agent":       {""},
				"post_raw":         {""},
				"custom_header":    {""},
				"host":             {""},
			},
		},
	}
	for resourceType, tt := range tests {
		t.Run(resourceType, func(t *testing.T) {
			api := newFakeAPI(t)
			config := func(body string) string {
				return providerBlock() + "\nresource \"" + resourceType + "\" \"test\" {\n" + body + "\n}"
			}

			runTest(t, api,
				resource.TestStep{Config: config(tt.full)},
				resource.TestStep{Config: config(tt.bare)},
			)

			var updates []apiRequest
			for _, r := range api.Requests() {
				if r.Method == "PUT" || r.Method == "PATCH" {
					updates = append(updates, r)
				}
			}
			if len(updates) != 1 {
				t.Fatalf("got %d updates, want 1", len(updates))
			}

			fields := make([]string, 0, len(tt.cleared))
			for f := range tt.cleared {
				fields = append(fields, f)
			}
			if diff := cmp.Diff(tt.cleared, pick(updates[0].Form, fields...)); diff != "" {
				t.Errorf("blanked settings differ (-want +got):\n%s", diff)
			}
		})
	}
}
