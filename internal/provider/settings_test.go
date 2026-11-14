//go:build dev

package provider_test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProviderSettingsFromTheEnvironment(t *testing.T) {
	api := newFakeAPI(t)

	// The first four requests are refused as too many, so the data source is only
	// read if the provider retries more often than it does by default, as it is
	// told to by the environment. It waits no longer than a second between
	// attempts.
	var failures atomic.Int32
	failures.Store(4)
	api.intercept = func(w http.ResponseWriter, r *http.Request) bool {
		if r.URL.Path != "/v1/uptime-locations" || failures.Load() == 0 {
			return false
		}
		failures.Add(-1)
		writeError(w, http.StatusTooManyRequests, "Too many requests", nil)
		return true
	}

	t.Setenv("STATUSCAKE_API_TOKEN", testToken)
	t.Setenv("STATUSCAKE_RPS", "1000")
	t.Setenv("STATUSCAKE_RETRIES", "5")
	t.Setenv("STATUSCAKE_MIN_BACKOFF", "0")
	t.Setenv("STATUSCAKE_MAX_BACKOFF", "1")

	runTest(t, api, resource.TestStep{
		Config: `
			provider "statuscake" {}

			data "statuscake_uptime_monitoring_locations" "test" {}
		`,
		Check: resource.TestCheckResourceAttr("data.statuscake_uptime_monitoring_locations.test", "locations.#", "2"),
	})

	calls := api.Calls("GET", "/v1/uptime-locations")
	if len(calls) < 5 {
		t.Fatalf("got %d requests, want the 4 that were refused and 1 more", len(calls))
	}
	if got, want := calls[0].Header.Get("Authorization"), "Bearer "+testToken; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}
