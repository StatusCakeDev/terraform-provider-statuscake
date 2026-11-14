//go:build dev

package provider_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestStateIsUnchanged guards what users' states hold. A resource is created
// from the smallest configuration there is, and the state that results holds
// every default and every value the API supplies. A state that differs from the
// one that users already have shows up as a change in their next plan.
//
// The state is that of the fake API, so the values the API supplies are the
// fake's. When a difference is deliberate, review it and refresh the files
// with:
//
//	go test -tags dev ./internal/provider -run TestStateIsUnchanged -update
func TestStateIsUnchanged(t *testing.T) {
	tests := map[string]struct {
		resource string
		config   string
	}{
		"contact_group": {
			resource: "statuscake_contact_group",
			config:   `name = "Operations"`,
		},
		"heartbeat_check": {
			resource: "statuscake_heartbeat_check",
			config: `
				name = "Job"
				period = 1800
			`,
		},
		"maintenance_window": {
			resource: "statuscake_maintenance_window",
			config: `
				name = "Window"
				start = "2022-01-29T00:00:00Z"
				end = "2022-01-30T23:59:59Z"
				timezone = "UTC"
				tags = ["production"]
			`,
		},
		"pagespeed_check": {
			resource: "statuscake_pagespeed_check",
			config: `
				name = "Page"
				check_interval = 300
				region = "UK"
				alert_config {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
		},
		"ssl_check": {
			resource: "statuscake_ssl_check",
			config: `
				check_interval = 600
				alert_config {
					alert_at = [7, 14, 21]
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
		},
		"uptime_check_http": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				http_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
		},
		"uptime_check_tcp": {
			resource: "statuscake_uptime_check",
			config: `
				name = "SSH"
				check_interval = 300
				tcp_check {
					port = 22
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
		},
		"uptime_check_dns": {
			resource: "statuscake_uptime_check",
			config: `
				name = "DNS"
				check_interval = 300
				dns_check {
					dns_ips = ["192.0.2.1"]
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
		},
		"uptime_check_icmp": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Ping"
				check_interval = 300
				icmp_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			api := newFakeAPI(t)
			addr := tt.resource + ".test"

			var attributes map[string]string
			runTest(t, api, resource.TestStep{
				Config: providerBlock() + "\nresource \"" + tt.resource + "\" \"test\" {\n  " + tt.config + "\n}",
				Check: func(s *terraform.State) error {
					attributes = s.RootModule().Resources[addr].Primary.Attributes
					return nil
				},
			})

			got, err := json.MarshalIndent(attributes, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, '\n')

			path := filepath.Join("testdata", "state", name+".golden.json")
			if *update {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, got, 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden file: %v (create it with -update)", err)
			}
			if diff := cmp.Diff(string(want), string(got)); diff != "" {
				t.Errorf("the state has changed (-golden +current):\n%s", diff)
			}
		})
	}
}
