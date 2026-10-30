//go:build dev

// TODO: fix this

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Some settings cannot be changed on a check that exists, and so changing them
// replaces it. Others are changed in place. Which is which is not in the schema
// that Terraform receives, but it decides whether a check, and the history it
// has, is deleted.
func TestUptimeCheckSettingsThatReplaceTheCheck(t *testing.T) {
	const (
		address = `monitored_resource { address = "https://example.com" }`
		base    = "name = \"Example\"\n  check_interval = 300\n  " + address + "\n  "
	)

	tests := map[string]struct {
		before string
		after  string
		want   plancheck.ResourceActionType
	}{
		"request method": {
			before: `
				http_check {
					request_method = "HTTP"
				}
			`,
			after: `
				http_check {
					request_method = "HEAD"
				}
			`,
			want: plancheck.ResourceActionReplace,
		},
		"tcp protocol": {
			before: `
				tcp_check {
					port = 22
					protocol = "TCP"
				}
			`,
			after: `
				tcp_check {
					port = 22
					protocol = "SSH"
			  }
			`,
			want: plancheck.ResourceActionReplace,
		},
		"http to icmp": {
			before: `http_check {}`,
			after:  `icmp_check {}`,
			want:   plancheck.ResourceActionReplace,
		},
		"icmp to dns": {
			before: `icmp_check {}`,
			after: `
				dns_check {
					dns_ips = ["192.0.2.1"]
				}
			`,
			want: plancheck.ResourceActionReplace,
		},
		"dns to tcp": {
			before: `
				dns_check {
					dns_ips = ["192.0.2.1"]
				}
			`,
			after: `
				tcp_check {
					port = 22
				}
			`,
			want: plancheck.ResourceActionReplace,
		},
		"http timeout": {
			before: `
				http_check {
					timeout = 10
				}
			`,
			after: `
				http_check {
					timeout = 20
				}
			`,
			want: plancheck.ResourceActionUpdate,
		},
		"tcp port": {
			before: `
				tcp_check {
					port = 22
				}
			`,
			after: `
				tcp_check {
					port = 25
				}
			`,
			want: plancheck.ResourceActionUpdate,
		},
		"dns ips": {
			before: `
				dns_check {
					dns_ips = ["192.0.2.1"]
				}
			`,
			after: `
				dns_check {
					dns_ips = ["192.0.2.2"]
				}
			`,
			want: plancheck.ResourceActionUpdate,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			api := newFakeAPI(t)
			config := func(check string) string {
				return providerBlock() + `resource "statuscake_uptime_check" "test" {
					` + base + check + `
				}`
			}

			runTest(t, api,
				resource.TestStep{Config: config(tt.before)},
				resource.TestStep{
					Config: config(tt.after),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("statuscake_uptime_check.test", tt.want),
						},
					},
				},
			)
		})
	}
}
