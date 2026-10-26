//go:build dev

package provider_test

import (
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These tests pin what the provider sends to the API. The API treats a field
// that is sent as one to change, so sending more than was configured, or less
// than has changed, alters what users' resources do.

func TestContactGroupSendsOnlyWhatIsConfigured(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_contact_group" "test" {
					name = "Operations"
				}
			`,
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_contact_group" "test" {
					name            = "Operations"
					email_addresses = ["ops@example.com"]
				}
			`,
		},
	)

	creates := api.Calls("POST", "/v1/contact-groups")
	if len(creates) != 1 {
		t.Fatalf("got %d creates, want 1", len(creates))
	}
	if diff := cmp.Diff(url.Values{"name": {"Operations"}}, creates[0].Form); diff != "" {
		t.Errorf("create body differs (-want +got):\n%s", diff)
	}

	var updates []apiRequest
	for _, r := range api.Requests() {
		if r.Method == "PUT" {
			updates = append(updates, r)
		}
	}
	if len(updates) != 1 {
		t.Fatalf("got %d updates, want 1", len(updates))
	}
	if diff := cmp.Diff(url.Values{"email_addresses[]": {"ops@example.com"}}, updates[0].Form); diff != "" {
		t.Errorf("update body differs (-want +got):\n%s", diff)
	}
}

func TestHeartbeatCheckClearsItsHostWhenTheBlockIsRemoved(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_heartbeat_check" "test" {
					name   = "Job"
					period = 1800

					monitored_resource {
						host = "gke"
					}
				}
			`,
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_heartbeat_check" "test" {
					name   = "Job"
					period = 1800
				}
			`,
		},
	)

	var updates []apiRequest
	for _, r := range api.Requests() {
		if r.Method == "PUT" {
			updates = append(updates, r)
		}
	}
	if len(updates) != 1 {
		t.Fatalf("got %d updates, want 1", len(updates))
	}
	if diff := cmp.Diff(url.Values{"host": {""}}, updates[0].Form); diff != "" {
		t.Errorf("update body differs (-want +got):\n%s", diff)
	}
}

func TestUptimeCheckCreateBody(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "Example"

					http_check {
						status_codes = ["200", "301"]
					}

					monitored_resource {
						address = "https://www.example.com"
					}

					regions = ["london"]
				}
			`,
		},
	)

	creates := api.Calls("POST", "/v1/uptime")
	if len(creates) != 1 {
		t.Fatalf("got %d creates, want 1", len(creates))
	}
	want := url.Values{
		"check_rate":       {"300"},
		"confirmation":     {"2"},
		"name":             {"Example"},
		"regions[]":        {"london"},
		"status_codes_csv": {"200,301"},
		"test_type":        {"HTTP"},
		"timeout":          {"15"},
		"trigger_rate":     {"0"},
		"website_url":      {"https://www.example.com"},
	}

	// The order of a set is not stable, so the status codes are compared as a
	// set.
	got := cloneForm(creates[0].Form)
	codes := strings.Split(got.Get("status_codes_csv"), ",")
	slices.Sort(codes)
	got.Set("status_codes_csv", strings.Join(codes, ","))

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("create body differs (-want +got):\n%s", diff)
	}
}
