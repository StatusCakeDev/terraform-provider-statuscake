//go:build dev

package provider_test

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPagespeedCheckLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_pagespeed_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_pagespeed_check" "test" {
					check_interval = 300
					name           = "Example"
					region         = "UK"

					alert_config {
						alert_bigger = "5000"
					}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "name", "Example"),
				resource.TestCheckResourceAttr(name, "check_interval", "300"),
				resource.TestCheckResourceAttr(name, "region", "UK"),
				resource.TestCheckResourceAttr(name, "alert_config.0.alert_bigger", "5000"),
				resource.TestCheckResourceAttr(name, "monitored_resource.0.address", "https://www.example.com"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_pagespeed_check" "test" {
					check_interval = 600
					name           = "Example"
					region         = "UK"
					paused         = true

					alert_config {
						alert_bigger = "5000"
					}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "check_interval", "600"),
				resource.TestCheckResourceAttr(name, "paused", "true"),
			),
		},
		// The region is not read back from the API, which reports a location
		// instead, so it is absent from an imported check.
		resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"region"}},
	)

	if got := api.Count("pagespeed"); got != 0 {
		t.Errorf("pagespeed checks left behind after destroy: %d", got)
	}
}

// Emptying the alert configuration resets every threshold, as an empty block
// disables all alerts.
func TestPagespeedCheckEmptyingAlertConfigResetsThresholds(t *testing.T) {
	api := newFakeAPI(t)
	config := func(alerts string) string {
		return providerBlock() + `
			resource "statuscake_pagespeed_check" "test" {
				check_interval = 300
				name           = "Example"
				region         = "UK"

				alert_config {` + alerts + `}

				monitored_resource {
					address = "https://www.example.com"
				}
			}
		`
	}

	runTest(t, api,
		resource.TestStep{Config: config(`
			alert_bigger  = 5000
			alert_slower  = 3000
			alert_smaller = 100
		`)},
		resource.TestStep{
			Config: config(""),
			Check:  resource.TestCheckResourceAttr("statuscake_pagespeed_check.test", "alert_config.0.alert_bigger", "0"),
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
	want := url.Values{"alert_bigger": {"0"}, "alert_slower": {"0"}, "alert_smaller": {"0"}}
	if diff := cmp.Diff(want, pick(updates[0].Form, "alert_bigger", "alert_slower", "alert_smaller")); diff != "" {
		t.Errorf("update body differs (-want +got):\n%s", diff)
	}
}
