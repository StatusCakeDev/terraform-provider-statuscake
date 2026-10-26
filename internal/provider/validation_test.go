package provider_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceConfigIsValidated(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		resource string
		config   string
		wantErr  *regexp.Regexp // nil when the config is valid
	}{
		"contact group: minimal": {
			resource: "statuscake_contact_group",
			config:   `name = "Ops"`,
			wantErr:  nil,
		},
		"contact group: name is required": {
			resource: "statuscake_contact_group",
			config:   "",
			wantErr:  rejected(`"name" is required`),
		},
		"contact group: bad email": {
			resource: "statuscake_contact_group",
			config: `
				name = "Ops"
				email_addresses = ["not-an-email"]
			`,
			wantErr: rejected("valid email address"),
		},
		"contact group: integration IDs are numbers": {
			resource: "statuscake_contact_group",
			config: `
				name = "Ops"
				integrations = ["abc"]
			`,
			wantErr: rejected("string representation of a number"),
		},
		"contact group: ping url must be http": {
			resource: "statuscake_contact_group",
			config: `
				name = "Ops"
			  ping_url = "ftp://example.com"
			`,
			wantErr: rejected("ping_url", "ftp"),
		},

		"heartbeat: minimal": {
			resource: "statuscake_heartbeat_check",
			config: `
				name = "Job"
			  period = 1800
			`,
			wantErr: nil,
		},
		"heartbeat: period too short": {
			resource: "statuscake_heartbeat_check",
			config: `
				name = "Job"
			  period = 29
			`,
			wantErr: rejected("expected period to be in the range"),
		},
		"heartbeat: period too long": {
			resource: "statuscake_heartbeat_check",
			config: `
				name = "Job"
			  period = 172801
			`,
			wantErr: rejected("expected period to be in the range"),
		},

		"uptime: http": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				http_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: nil,
		},
		"uptime: a check type is required": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("one of", "dns_check", "http_check"),
		},
		"uptime: only one check type": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				http_check {}
				icmp_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("only one of"),
		},
		"uptime: unsupported interval": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 7
				http_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("check_interval"),
		},
		"uptime: timeout bounds": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				http_check {
					timeout = 4
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("expected http_check.0.timeout to be in the range"),
		},
		"uptime: confirmation bounds": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				confirmation = 4
				http_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("expected confirmation to be in the range"),
		},
		"uptime: tcp needs a port": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				tcp_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected(`"port" is required`),
		},
		"uptime: tcp protocol": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				tcp_check {
					port = 22
					protocol = "FTP"
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("expected tcp_check.0.protocol to be one of"),
		},
		"uptime: dns needs addresses": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				dns_check {}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected(`"dns_ips" is required`),
		},
		"uptime: dns addresses are IPs": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				dns_check {
					dns_ips = ["example.com"]
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("to contain a valid IP"),
		},
		"uptime: payload and raw payload conflict": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				http_check {
					request_payload = {
						a = "b"
					}
					request_payload_raw = "a=b"
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("conflicts with"),
		},
		"uptime: matcher is one of two": {
			resource: "statuscake_uptime_check",
			config: `
				name = "Example"
				check_interval = 300
				http_check {
					content_matchers {
						content = "x"
						matcher = "MAYBE"
					}
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: rejected("expected http_check.0.content_matchers.0.matcher to be one of"),
		},

		"ssl: minimal": {
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
			wantErr: nil,
		},
		"ssl: needs an address": {
			resource: "statuscake_ssl_check",
			config: `
				check_interval = 600
				alert_config {
					alert_at = [7, 14, 21]
				}
			`,
			wantErr: rejected("monitored_resource", "required")},

		"pagespeed: minimal": {
			resource: "statuscake_pagespeed_check",
			config: `
				name = "Page"
			  check_interval = 300
			  region = "UK"
			  alert_config {
					alert_bigger = 5000
				}
				monitored_resource {
					address = "https://example.com"
				}
			`,
			wantErr: nil,
		},
		"pagespeed: unknown region": {
			resource: "statuscake_pagespeed_check",
			config: `
				name = "Page"
			  check_interval = 300
			  region = "Atlantis"
			  alert_config {
					alert_bigger = 5000
				}
				monitored_resource {
					address = "https://example.com"
				}
		  `,
			wantErr: rejected("expected region to be one of"),
		},

		"maintenance window: needs tags or tests": {
			resource: "statuscake_maintenance_window",
			config: `
				name = "W"
			  start = "2022-01-29T00:00:00Z"
			  end = "2022-01-30T00:00:00Z"
			  timezone = "UTC"
			`,
			wantErr: rejected("one of", "tags", "tests"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			step := resource.TestStep{
				Config:   providerBlock() + "\nresource \"" + tt.resource + "\" \"test\" {\n  " + tt.config + "\n}",
				PlanOnly: true,
			}
			if tt.wantErr == nil {
				step.ExpectNonEmptyPlan = true
			} else {
				step.ExpectError = tt.wantErr
			}
			runSteps(t, step)
		})
	}
}

func TestProviderConfigIsValidated(t *testing.T) {
	t.Parallel()

	// Planning a data source is enough to configure the provider, and the
	// configuration is rejected before the data source is read.
	const data = `data "statuscake_uptime_monitoring_locations" "test" {}`
	token := `api_token = "` + testToken + `"`

	tests := map[string]struct {
		settings string
		wantErr  *regexp.Regexp
	}{
		"token is required": {
			settings: "",
			wantErr:  rejected(`"api_token" is required`),
		},
		"token characters": {
			settings: `api_token = "has spaces and is long enough"`,
			wantErr:  rejected("API token must only contain"),
		},
		"rps at least one": {
			settings: token + `
				rps = 0
			`,
			wantErr: rejected("expected rps to be at least (1)"),
		},
		"retries at most ten": {
			settings: token + `
				retries = 11
			`,
			wantErr: rejected("expected retries to be in the range (0 - 10)"),
		},
		"min backoff not negative": {
			settings: token + `
				min_backoff = -1
			`,
			wantErr: rejected("expected min_backoff to be at least (0)"),
		},
		"max backoff at least one": {
			settings: token + `
				max_backoff = 0
			`,
			wantErr: rejected("expected max_backoff to be at least (1)"),
		},
		"no custom endpoint in a release": {
			settings: token + `
				custom_endpoint = "https://api.example.com"
			`,
			wantErr: rejected("custom_endpoint"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if name == "no custom endpoint in a release" && developmentBuild {
				t.Skip("this is a development build, which has the setting")
			}

			runSteps(t, resource.TestStep{
				Config:      "provider \"statuscake\" {\n  " + tt.settings + "\n}\n" + data,
				PlanOnly:    true,
				ExpectError: tt.wantErr,
			})
		})
	}
}

// rejected matches an error that mentions each of the words, whichever way
// Terraform wraps the lines of its message.
func rejected(words ...string) *regexp.Regexp {
	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = strings.ReplaceAll(regexp.QuoteMeta(w), " ", `[\s│]+`)
	}
	return regexp.MustCompile("(?s)" + strings.Join(parts, ".*"))
}
