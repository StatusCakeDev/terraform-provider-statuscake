//go:build dev

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMaintenanceWindowLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_maintenance_window.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_maintenance_window" "test" {
					name            = "Weekends"
					start           = "2022-01-29T00:00:00Z"
					end             = "2022-01-30T23:59:59Z"
					repeat_interval = "1w"
					timezone        = "UTC"
					tags            = ["production"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "name", "Weekends"),
				resource.TestCheckResourceAttr(name, "repeat_interval", "1w"),
				resource.TestCheckResourceAttr(name, "timezone", "UTC"),
				resource.TestCheckResourceAttr(name, "tags.#", "1"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_maintenance_window" "test" {
					name     = "Weekends"
					start    = "2022-01-29T00:00:00Z"
					end      = "2022-01-30T23:59:59Z"
					timezone = "Europe/London"
					tags     = ["staging"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "repeat_interval", "never"),
				resource.TestCheckResourceAttr(name, "timezone", "Europe/London"),
				resource.TestCheckResourceAttr(name, "tags.0", "staging"),
			),
		},
		resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true},
	)

	if got := api.Count("maintenance-windows"); got != 0 {
		t.Errorf("maintenance windows left behind after destroy: %d", got)
	}
}
