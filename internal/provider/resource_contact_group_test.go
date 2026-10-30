//go:build dev

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestContactGroupLifecycle(t *testing.T) {
	// Not parallel: the tests share the process environment that the testing
	// framework reads.
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_contact_group" "test" {
					name            = "Operations"
					email_addresses = ["ops@example.com", "oncall@example.com"]
					mobile_numbers  = ["+447700900123"]
					ping_url        = "https://example.com/hook"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("statuscake_contact_group.test", "id"),
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "name", "Operations"),
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "email_addresses.#", "2"),
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "mobile_numbers.#", "1"),
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "ping_url", "https://example.com/hook"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_contact_group" "test" {
					name            = "Operations (renamed)"
					email_addresses = ["ops@example.com"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "name", "Operations (renamed)"),
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "email_addresses.#", "1"),
				resource.TestCheckResourceAttr("statuscake_contact_group.test", "mobile_numbers.#", "0"),
			),
		},
		resource.TestStep{
			ResourceName:      "statuscake_contact_group.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	)

	if got := api.Count("contact-groups"); got != 0 {
		t.Errorf("contact groups left behind after destroy: %d", got)
	}
}
