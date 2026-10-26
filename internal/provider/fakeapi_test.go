//go:build dev

package provider_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// apiRequest is a request that the fake received.
type apiRequest struct {
	Method string
	Path   string
	Query  url.Values
	Form   url.Values
	Header http.Header
}

// fakeAPI stands in for the StatusCake API. It keeps the resources that are
// created in memory and answers the way the real API does, so that the provider
// can be exercised end to end without creating anything real.
//
// The fake is deliberately simple: it stores the form that was sent and renders
// it back as the API would. It does not validate beyond what a test asks of it.
type fakeAPI struct {
	URL string

	mu       sync.Mutex
	seq      int
	items    map[string]map[string]url.Values
	requests []apiRequest

	// intercept, when set, may answer a request itself. It returns true when it
	// has done so.
	intercept func(w http.ResponseWriter, r *http.Request) bool
}

// newFakeAPI starts a fake that is stopped when the test ends.
func newFakeAPI(t *testing.T) *fakeAPI {
	t.Helper()

	f := &fakeAPI{items: map[string]map[string]url.Values{}}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)
	f.URL = srv.URL
	return f
}

// Requests returns the requests received so far, in order.
func (f *fakeAPI) Requests() []apiRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.requests)
}

// Calls returns the requests received that match the method and path.
func (f *fakeAPI) Calls(method, path string) []apiRequest {
	var out []apiRequest
	for _, r := range f.Requests() {
		if r.Method == method && r.Path == path {
			out = append(out, r)
		}
	}
	return out
}

// Count returns how many resources of the kind exist, such as "uptime".
func (f *fakeAPI) Count(kind string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.items[kind])
}

// Clear removes every resource of the kind, as if they had been deleted outside
// of Terraform.
func (f *fakeAPI) Clear(kind string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.items[kind] = map[string]url.Values{}
}

// Seed adds a resource that already exists, returning its ID.
func (f *fakeAPI) Seed(kind string, form url.Values) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.insert(kind, form)
}

func (f *fakeAPI) insert(kind string, form url.Values) string {
	f.seq++
	id := strconv.Itoa(1000 + f.seq)
	if f.items[kind] == nil {
		f.items[kind] = map[string]url.Values{}
	}
	f.items[kind][id] = form
	return id
}

func (f *fakeAPI) serve(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "bad form", nil)
		return
	}

	f.mu.Lock()
	f.requests = append(f.requests, apiRequest{Method: r.Method, Path: r.URL.Path, Query: cloneForm(r.URL.Query()), Form: cloneForm(r.PostForm), Header: r.Header.Clone()})
	intercept := f.intercept
	f.mu.Unlock()

	if intercept != nil && intercept(w, r) {
		return
	}

	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1"), "/"), "/")
	kind := parts[0]
	id := ""
	if len(parts) > 1 {
		id = parts[1]
	}

	switch {
	case r.Method == http.MethodGet && kind == "uptime-locations", r.Method == http.MethodGet && kind == "pagespeed-locations":
		writeJSON(w, http.StatusOK, map[string]any{"data": locations(kind)})
	case r.Method == http.MethodPost && id == "":
		f.mu.Lock()
		newID := f.insert(kind, cloneForm(r.PostForm))
		f.mu.Unlock()
		writeJSON(w, http.StatusCreated, map[string]any{"data": map[string]any{"new_id": newID}})
	case r.Method == http.MethodGet && id != "":
		f.mu.Lock()
		form, ok := f.items[kind][id]
		f.mu.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "Not Found", nil)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"data": render(kind, id, form)})
	case r.Method == http.MethodPut && id != "":
		f.mu.Lock()
		form, ok := f.items[kind][id]
		if ok {
			merge(form, r.PostForm)
		}
		f.mu.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "Not Found", nil)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete && id != "":
		f.mu.Lock()
		_, ok := f.items[kind][id]
		delete(f.items[kind], id)
		f.mu.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "Not Found", nil)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusNotFound, "no such route", nil)
	}
}

// merge applies an update: fields that are sent replace those that are stored,
// and a list replaces the whole list. A list holding a single empty value, which
// is how the client says "make this empty", clears it.
func merge(stored, update url.Values) {
	for k, v := range update {
		if strings.HasSuffix(k, "[]") && len(v) == 1 && v[0] == "" {
			delete(stored, k)
			continue
		}
		stored[k] = slices.Clone(v)
	}
}

func cloneForm(v url.Values) url.Values {
	out := make(url.Values, len(v))
	for k, vs := range v {
		out[k] = slices.Clone(vs)
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string, errs map[string][]string) {
	body := map[string]any{"message": message}
	if errs != nil {
		body["errors"] = errs
	}
	writeJSON(w, status, body)
}

// Helpers that read a stored form the way the API would.

func str(f url.Values, k string) string { return f.Get(k) }

func optStr(f url.Values, k string) any {
	if v := f.Get(k); v != "" {
		return v
	}
	return nil
}

func list(f url.Values, k string) []string {
	out := []string{}
	for _, v := range f[k+"[]"] {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func boolean(f url.Values, k string) bool {
	b, _ := strconv.ParseBool(f.Get(k))
	return b
}

func integer(f url.Values, k string) int {
	n, _ := strconv.Atoi(f.Get(k))
	return n
}

func optInt(f url.Values, k string) any {
	if f.Get(k) == "" {
		return nil
	}
	return integer(f, k)
}

func ints(f url.Values, k string) []int {
	out := []int{}
	for _, v := range list(f, k) {
		n, _ := strconv.Atoi(v)
		out = append(out, n)
	}
	return out
}

// render turns a stored form into the response the API gives for it.
func render(kind, id string, f url.Values) map[string]any {
	switch kind {
	case "contact-groups":
		return map[string]any{
			"id":              id,
			"name":            str(f, "name"),
			"email_addresses": list(f, "email_addresses"),
			"integrations":    list(f, "integrations"),
			"mobile_numbers":  list(f, "mobile_numbers"),
			"ping_url":        optStr(f, "ping_url"),
		}
	case "heartbeat":
		return map[string]any{
			"id":             id,
			"name":           str(f, "name"),
			"url":            "https://heartbeat.example.com/" + id,
			"period":         integer(f, "period"),
			"contact_groups": list(f, "contact_groups"),
			"host":           optStr(f, "host"),
			"paused":         boolean(f, "paused"),
			"status":         "up",
			"tags":           list(f, "tags"),
			"uptime":         100,
		}
	case "maintenance-windows":
		return map[string]any{
			"id":              id,
			"name":            str(f, "name"),
			"start_at":        str(f, "start_at"),
			"end_at":          str(f, "end_at"),
			"repeat_interval": str(f, "repeat_interval"),
			"state":           "pending",
			"tags":            list(f, "tags"),
			"tests":           list(f, "tests"),
			"timezone":        str(f, "timezone"),
		}
	case "pagespeed":
		return map[string]any{
			"id":             id,
			"name":           str(f, "name"),
			"website_url":    str(f, "website_url"),
			"check_rate":     integer(f, "check_rate"),
			"alert_bigger":   integer(f, "alert_bigger"),
			"alert_slower":   integer(f, "alert_slower"),
			"alert_smaller":  integer(f, "alert_smaller"),
			"contact_groups": list(f, "contact_groups"),
			"location":       str(f, "region"),
			"paused":         boolean(f, "paused"),
		}
	case "ssl":
		return map[string]any{
			"id":               id,
			"website_url":      str(f, "website_url"),
			"check_rate":       integer(f, "check_rate"),
			"alert_at":         ints(f, "alert_at"),
			"alert_broken":     boolean(f, "alert_broken"),
			"alert_expiry":     boolean(f, "alert_expiry"),
			"alert_mixed":      boolean(f, "alert_mixed"),
			"alert_reminder":   boolean(f, "alert_reminder"),
			"contact_groups":   list(f, "contact_groups"),
			"follow_redirects": boolean(f, "follow_redirects"),
			"hostname":         optStr(f, "hostname"),
			"mixed_content":    []any{},
			"paused":           boolean(f, "paused"),
			"user_agent":       optStr(f, "user_agent"),
			"flags":            map[string]any{},
		}
	case "uptime":
		servers := []map[string]any{}
		for _, code := range list(f, "regions") {
			servers = append(servers, map[string]any{"region": code, "region_code": code, "description": code, "status": "up"})
		}
		codes := []string{}
		for _, c := range strings.Split(f.Get("status_codes_csv"), ",") {
			if c = strings.TrimSpace(c); c != "" {
				codes = append(codes, c)
			}
		}
		return map[string]any{
			"id":               id,
			"name":             str(f, "name"),
			"test_type":        str(f, "test_type"),
			"website_url":      str(f, "website_url"),
			"check_rate":       integer(f, "check_rate"),
			"confirmation":     integer(f, "confirmation"),
			"contact_groups":   list(f, "contact_groups"),
			"custom_header":    optStr(f, "custom_header"),
			"dns_ips":          list(f, "dns_ips"),
			"dns_server":       optStr(f, "dns_server"),
			"do_not_find":      boolean(f, "do_not_find"),
			"enable_ssl_alert": boolean(f, "enable_ssl_alert"),
			"final_endpoint":   optStr(f, "final_endpoint"),
			"find_string":      optStr(f, "find_string"),
			"follow_redirects": boolean(f, "follow_redirects"),
			"include_header":   boolean(f, "include_header"),
			"host":             optStr(f, "host"),
			"paused":           boolean(f, "paused"),
			"port":             optInt(f, "port"),
			"post_body":        optStr(f, "post_body"),
			"post_raw":         optStr(f, "post_raw"),
			"processing":       false,
			"servers":          servers,
			"status":           "up",
			"status_codes":     codes,
			"tags":             list(f, "tags"),
			"timeout":          integer(f, "timeout"),
			"trigger_rate":     integer(f, "trigger_rate"),
			"uptime":           100,
			"use_jar":          boolean(f, "use_jar"),
			"user_agent":       optStr(f, "user_agent"),
		}
	}
	return map[string]any{"id": id}
}

func locations(kind string) []map[string]any {
	return []map[string]any{
		{"description": "London, UK", "region": "london", "region_code": "london", "status": "up", "ipv4": "192.0.2.1"},
		{"description": "New York, US", "region": "new-york", "region_code": "new-york", "status": "up", "ipv4": "192.0.2.2", "ipv6": "2001:db8::2"},
	}
}
