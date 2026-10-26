//go:build dev

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSSLCheckLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_ssl_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_ssl_check" "test" {
					check_interval = 600
					user_agent     = "terraform managed SSL check"

					alert_config {
						alert_at    = [7, 14, 21]
						on_broken   = false
						on_expiry   = true
						on_mixed    = false
						on_reminder = true
					}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "check_interval", "600"),
				resource.TestCheckResourceAttr(name, "alert_config.0.alert_at.#", "3"),
				resource.TestCheckResourceAttr(name, "alert_config.0.on_expiry", "true"),
				resource.TestCheckResourceAttr(name, "alert_config.0.on_reminder", "true"),
				resource.TestCheckResourceAttr(name, "monitored_resource.0.address", "https://www.example.com"),
				resource.TestCheckResourceAttr(name, "user_agent", "terraform managed SSL check"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_ssl_check" "test" {
					check_interval = 1800

					alert_config {
						alert_at   = [1, 7, 30]
						on_expiry  = true
					}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "check_interval", "1800"),
				resource.TestCheckResourceAttr(name, "alert_config.0.alert_at.#", "3"),
			),
		},
		resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true},
	)

	if got := api.Count("ssl"); got != 0 {
		t.Errorf("SSL checks left behind after destroy: %d", got)
	}
}
