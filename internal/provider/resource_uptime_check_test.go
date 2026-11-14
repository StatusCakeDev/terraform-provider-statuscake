//go:build dev

package provider_test

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUptimeCheckHTTPLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_uptime_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					confirmation   = 3
					name           = "Example"
					trigger_rate   = 10

					http_check {
						enable_cookies   = true
						follow_redirects = true
						timeout          = 20
						user_agent       = "terraform managed uptime check"
						validate_ssl     = true

						basic_authentication {
							password = "password"
							username = "username"
						}

						content_matchers {
							content         = "Welcome"
							include_headers = true
							matcher         = "CONTAINS_STRING"
						}

						request_headers = {
							Authorization = "bearer 123456"
						}

						status_codes = ["202", "404", "405"]
					}

					monitored_resource {
						address = "https://www.example.com"
					}

					regions = ["london", "paris"]
					tags    = ["production"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "name", "Example"),
				resource.TestCheckResourceAttr(name, "check_interval", "300"),
				resource.TestCheckResourceAttr(name, "confirmation", "3"),
				resource.TestCheckResourceAttr(name, "trigger_rate", "10"),
				resource.TestCheckResourceAttr(name, "monitored_resource.0.address", "https://www.example.com"),
				resource.TestCheckResourceAttr(name, "http_check.0.timeout", "20"),
				resource.TestCheckResourceAttr(name, "http_check.0.status_codes.#", "3"),
				resource.TestCheckResourceAttr(name, "http_check.0.content_matchers.0.content", "Welcome"),
				resource.TestCheckResourceAttr(name, "http_check.0.request_headers.Authorization", "bearer 123456"),
				resource.TestCheckResourceAttr(name, "regions.#", "2"),
				resource.TestCheckResourceAttr(name, "tags.#", "1"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 60
					confirmation   = 2
					name           = "Example (renamed)"
					paused         = true

					http_check {
						follow_redirects = false
						timeout          = 15
						status_codes     = ["200"]

						content_matchers {
							content = "Welcome"
						}
					}

					monitored_resource {
						address = "https://www.example.com"
					}

					regions = ["london"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "name", "Example (renamed)"),
				resource.TestCheckResourceAttr(name, "check_interval", "60"),
				resource.TestCheckResourceAttr(name, "paused", "true"),
				resource.TestCheckResourceAttr(name, "http_check.0.status_codes.#", "1"),
				resource.TestCheckResourceAttr(name, "regions.#", "1"),
				resource.TestCheckResourceAttr(name, "tags.#", "0"),
			),
		},
	)

	if got := api.Count("uptime"); got != 0 {
		t.Errorf("uptime checks left behind after destroy: %d", got)
	}
}

func TestUptimeCheckTCPLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_uptime_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "SSH"

					tcp_check {
						port     = 22
						protocol = "SSH"
						timeout  = 10
					}

					monitored_resource {
						address = "bastion.example.com"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "tcp_check.0.port", "22"),
				resource.TestCheckResourceAttr(name, "tcp_check.0.protocol", "SSH"),
				resource.TestCheckResourceAttr(name, "tcp_check.0.timeout", "10"),
				resource.TestCheckResourceAttr(name, "monitored_resource.0.address", "bastion.example.com"),
			),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "SSH"

					tcp_check {
						port     = 22
						protocol = "SSH"
						timeout  = 30
					}

					monitored_resource {
						address = "bastion.example.com"
					}
				}
			`,
			Check: resource.TestCheckResourceAttr(name, "tcp_check.0.timeout", "30"),
		},
	)
}

func TestUptimeCheckDNSLifecycle(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_uptime_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "DNS"

					dns_check {
						dns_ips    = ["192.0.2.10", "192.0.2.11"]
						dns_server = "8.8.8.8"
					}

					monitored_resource {
						address = "www.example.com"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "dns_check.0.dns_ips.#", "2"),
				resource.TestCheckResourceAttr(name, "dns_check.0.dns_server", "8.8.8.8"),
			),
		},
	)
}

func TestUptimeCheckICMPLifecycle(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "Ping"

					icmp_check {}

					monitored_resource {
						address = "192.0.2.1"
					}
				}
			`,
			Check: resource.TestCheckResourceAttr("statuscake_uptime_check.test", "monitored_resource.0.address", "192.0.2.1"),
		},
	)
}

// Removing the content matchers from a check clears them from the API. Nothing
// used to be sent for an empty block, so the check kept the matcher and every
// later plan tried to remove it again.
func TestUptimeCheckRemovingContentMatchersClearsThem(t *testing.T) {
	api := newFakeAPI(t)
	const name = "statuscake_uptime_check.test"

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "Example"

					http_check {
						content_matchers {
							content         = "Welcome"
							include_headers = true
							matcher         = "NOT_CONTAINS_STRING"
						}
					}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
			Check: resource.TestCheckResourceAttr(name, "http_check.0.content_matchers.#", "1"),
		},
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "Example"

					http_check {}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
			Check: resource.TestCheckResourceAttr(name, "http_check.0.content_matchers.#", "0"),
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
	want := url.Values{
		"find_string":    {""},
		"include_header": {"false"},
		"do_not_find":    {"false"},
	}
	if diff := cmp.Diff(want, updates[0].Form); diff != "" {
		t.Errorf("update body differs (-want +got):\n%s", diff)
	}
}

// Leaving the block out when the check is created sends nothing about it.
func TestUptimeCheckWithoutContentMatchersSendsNoMatcher(t *testing.T) {
	api := newFakeAPI(t)

	runTest(t, api,
		resource.TestStep{
			Config: providerBlock() + `
				resource "statuscake_uptime_check" "test" {
					check_interval = 300
					name           = "Example"

					http_check {}

					monitored_resource {
						address = "https://www.example.com"
					}
				}
			`,
		},
	)

	creates := api.Calls("POST", "/v1/uptime")
	if len(creates) != 1 {
		t.Fatalf("got %d creates, want 1", len(creates))
	}
	for _, field := range []string{"find_string", "include_header", "do_not_find"} {
		if _, ok := creates[0].Form[field]; ok {
			t.Errorf("%s was sent when no matcher was configured", field)
		}
	}
}

// Removing credentials from a check blanks them on the API. Nothing used to be
// sent for the removed block, so the credentials stayed on the check.
func TestUptimeCheckRemovingBasicAuthenticationClearsIt(t *testing.T) {
	tests := map[string]struct {
		with    string
		without string
	}{
		"http": {
			with: `
				http_check {
					basic_authentication {
						username = "admin"
						password = "secret"
					}
				}
			`,
			without: `
				http_check {}
			`,
		},
		"tcp": {
			with: `
				tcp_check {
					port = 25
					authentication {
						username = "admin"
						password = "secret"
					}
				}
			`,
			without: `
				tcp_check {
					port = 25
				}
			`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			api := newFakeAPI(t)
			config := func(check string) string {
				return providerBlock() + `
					resource "statuscake_uptime_check" "test" {
						check_interval = 300
						name           = "Example"
						` + check + `

						monitored_resource {
							address = "example.com"
						}
					}
				`
			}

			runTest(t, api,
				resource.TestStep{Config: config(tt.with)},
				resource.TestStep{Config: config(tt.without)},
			)

			creates := api.Calls("POST", "/v1/uptime")
			if len(creates) != 1 {
				t.Fatalf("got %d creates, want 1", len(creates))
			}
			if diff := cmp.Diff(url.Values{"basic_username": {"admin"}, "basic_password": {"secret"}},
				pick(creates[0].Form, "basic_username", "basic_password")); diff != "" {
				t.Errorf("credentials on create differ (-want +got):\n%s", diff)
			}

			var updates []apiRequest
			for _, r := range api.Requests() {
				if r.Method == "PUT" {
					updates = append(updates, r)
				}
			}
			if len(updates) != 1 {
				t.Fatalf("got %d updates, want 1", len(updates))
			}
			want := url.Values{"basic_username": {""}, "basic_password": {""}}
			if diff := cmp.Diff(want, pick(updates[0].Form, "basic_username", "basic_password")); diff != "" {
				t.Errorf("credentials on update differ (-want +got):\n%s", diff)
			}
		})
	}
}

// pick returns the given fields of a form.
func pick(form url.Values, fields ...string) url.Values {
	out := url.Values{}
	for _, f := range fields {
		if v, ok := form[f]; ok {
			out[f] = v
		}
	}
	return out
}
