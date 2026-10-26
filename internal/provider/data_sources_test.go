//go:build dev

package provider_test

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestContactGroupDataSource(t *testing.T) {
	api := newFakeAPI(t)
	id := api.Seed("contact-groups", url.Values{
		"name":              {"Operations"},
		"email_addresses[]": {"ops@example.com"},
		"ping_url":          {"https://example.com/hook"},
	})

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				data "statuscake_contact_group" "test" {
					id = "` + id + `"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.statuscake_contact_group.test", "name", "Operations"),
				resource.TestCheckResourceAttr("data.statuscake_contact_group.test", "email_addresses.#", "1"),
				resource.TestCheckResourceAttr("data.statuscake_contact_group.test", "ping_url", "https://example.com/hook"),
			),
		},
	)
}

func TestMonitoringLocationsDataSources(t *testing.T) {
	tests := map[string]struct {
		path  string
		param string
	}{
		"statuscake_uptime_monitoring_locations": {
			path:  "/v1/uptime-locations",
			param: "region_code",
		},
		"statuscake_pagespeed_monitoring_locations": {
			path:  "/v1/pagespeed-locations",
			param: "location",
		},
	}
	for dataSource, tt := range tests {
		t.Run(dataSource, func(t *testing.T) {
			api := newFakeAPI(t)
			addr := "data." + dataSource + ".test"

			runTest(t, api,
				resource.TestStep{
					Config: providerBlock() + `
						data "` + dataSource + `" "test" {
							region_code = "GBR"
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(addr, "locations.#", "2"),
						resource.TestCheckResourceAttr(addr, "locations.0.region_code", "london"),
						resource.TestCheckResourceAttr(addr, "locations.1.ipv6", "2001:db8::2"),
					),
				},
			)

			var got []url.Values
			for _, r := range api.Calls("GET", tt.path) {
				got = append(got, r.Query)
			}

			if len(got) == 0 {
				t.Fatalf("no request was made to %s", tt.path)
			}

			if diff := cmp.Diff(url.Values{tt.param: {"GBR"}}, got[0]); diff != "" {
				t.Errorf("query sent to the API differs (-want +got):\n%s", diff)
			}
		})
	}
}
