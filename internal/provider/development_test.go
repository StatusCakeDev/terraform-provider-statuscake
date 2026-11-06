//go:build dev

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These tests cover the settings that only development builds have, and so run
// only with `-tags dev`.

func TestDevelopmentCustomEndpointIsValidated(t *testing.T) {
	t.Parallel()

	for endpoint, valid := range map[string]bool{
		"https://api.example.com": true,
		"http://localhost:8080":   true,
		"localhost":               false,
		"ftp://example.com":       false,
	} {
		t.Run(endpoint, func(t *testing.T) {
			t.Parallel()

			// A resource is planned without reaching for the endpoint, which a
			// valid one has nothing behind it to answer.
			step := resource.TestStep{
				Config: fmt.Sprintf(`
					provider "statuscake" {
						api_token       = %q
						custom_endpoint = %q
					}

					resource "statuscake_contact_group" "test" {
						name = "Operations"
					}
				`, testToken, endpoint),
				PlanOnly: true,
			}
			if valid {
				step.ExpectNonEmptyPlan = true
			} else {
				step.ExpectError = rejected("custom_endpoint")
			}
			runSteps(t, step)
		})
	}
}

func TestDevelopmentCustomEndpointIsUsed(t *testing.T) {
	api := newFakeAPI(t)
	requireCLI(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
				provider "statuscake" {
					api_token       = %q
					custom_endpoint = %q
					rps             = 1000
					retries         = 0
				}

				data "statuscake_uptime_monitoring_locations" "test" {}
			`, testToken, api.URL),
			Check: resource.TestCheckResourceAttr("data.statuscake_uptime_monitoring_locations.test", "locations.#", "2"),
		}},
	})

	if len(api.Calls("GET", "/v1/uptime-locations")) == 0 {
		t.Error("the request did not reach the endpoint that was configured")
	}
}

func TestDevelopmentCustomEndpointFromTheEnvironment(t *testing.T) {
	api := newFakeAPI(t)

	// runTest points the provider at the fake through the environment.
	runTest(t, api, resource.TestStep{
		Config: providerBlock() + `data "statuscake_uptime_monitoring_locations" "test" {}`,
		Check:  resource.TestCheckResourceAttr("data.statuscake_uptime_monitoring_locations.test", "locations.#", "2"),
	})

	if len(api.Calls("GET", "/v1/uptime-locations")) == 0 {
		t.Error("the request did not reach the endpoint in the environment")
	}
}

// runTest runs the steps against a fake API, to which the provider is pointed
// by the setting that development builds have. No real resources are created.
func runTest(t *testing.T, api *fakeAPI, steps ...resource.TestStep) {
	t.Helper()
	requireCLI(t)

	t.Setenv("STATUSCAKE_CUSTOM_ENDPOINT", api.URL)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: providerFactories(),
		Steps:                    steps,
	})
}
