//go:build dev

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestHeartbeatCheckLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_heartbeat_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_heartbeat_check" "test" {
					name   = "Nightly job"
					period = 1800
					tags   = ["production"]

					monitored_resource {
						host = "gke"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "name", "Nightly job"),
				resource.TestCheckResourceAttr(name, "period", "1800"),
				resource.TestCheckResourceAttr(name, "paused", "false"),
				resource.TestCheckResourceAttr(name, "tags.#", "1"),
				resource.TestCheckResourceAttr(name, "monitored_resource.0.host", "gke"),
				resource.TestCheckResourceAttrSet(name, "check_url"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_heartbeat_check" "test" {
					name   = "Nightly job"
					period = 3600
					paused = true
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "period", "3600"),
				resource.TestCheckResourceAttr(name, "paused", "true"),
				resource.TestCheckResourceAttr(name, "tags.#", "0"),
				resource.TestCheckResourceAttr(name, "monitored_resource.#", "0"),
			),
		},
		resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true},
	)

	if got := api.Count("heartbeat"); got != 0 {
		t.Errorf("heartbeat checks left behind after destroy: %d", got)
	}
}
