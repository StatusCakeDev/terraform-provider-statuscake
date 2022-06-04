---
layout: "statuscake"
page_title: "StatusCake: statuscake_test"
sidebar_current: "docs-statuscake-test"
description: |-
  The statuscake_test resource allows StatusCake tests to be managed by Terraform.
---

# statuscake\_test

The test resource allows StatusCake tests to be managed by Terraform.

## Example Usage

```hcl
resource "statuscake_test" "google" {
  website_name  = "google.com"
  website_url   = "www.google.com"
  test_type     = "HTTP"
  check_rate    = 300
  contact_group = ["12345"]
}
```

## Argument Reference

The following arguments are supported:

* `website_name` - (Required) This is the name of the test and the website to be monitored.
* `website_url` - (Required) The URL of the website to be monitored
* `check_rate` - (Optional) Test check rate in seconds. Defaults to 300
* `contact_id` - **Deprecated** (Optional) The id of the contact group to be added to the test. Each test can have only one.
* `contact_group` - (Optional) Set test contact groups, must be array of strings.
* `test_type` - (Required) The type of Test. Either HTTP, TCP, PING, or DNS
* `paused` - (Optional) Whether or not the test is paused. Defaults to false.
* `timeout` - (Optional) The timeout of the test in seconds.
* `confirmations` - (Optional) The number of confirmation servers to use in order to detect downtime. Defaults to 0.
* `port` - (Optional) The port to use when specifying a TCP test.
* `trigger_rate` - (Optional) The number of minutes to wait before sending an alert. Default is `5`.
* `custom_header` - (Optional) Custom HTTP header, must be supplied as JSON.
* `user_agent` - (Optional) Test with a custom user agent set.
* `node_locations` - (Optional) Set test node locations, must be array of strings.
* `ping_url` - (Optional) A URL to ping if a site goes down.
* `basic_user` - (Optional) A Basic Auth User account to use to login
* `basic_pass` - (Optional) If BasicUser is set then this should be the password for the BasicUser.
* `public` - (Optional) Set 1 to enable public reporting, 0 to disable.
* `logo_image` - (Optional) A URL to a image to use for public reporting.
* `branding` - (Optional) Set to 0 to use branding (default) or 1 to disable public reporting branding).
* `website_host` - (Optional) Used internally, when possible please add.
* `virus` - (Optional) Enable virus checking or not. 1 to enable
* `find_string` - (Optional) A string that should either be found or not found.
* `do_not_find` - (Optional) If the above string should be found to trigger a alert. 1 = will trigger if find_string found.
* `real_browser` - (Optional) Use 1 to TURN OFF real browser testing.
* `test_tags` - (Optional) Set test tags, must be array of strings.
* `status_codes` - (Optional) Comma Separated List of StatusCodes to Trigger Error on. Defaults are "204, 205, 206, 303, 400, 401, 403, 404, 405, 406, 408, 410, 413, 444, 429, 494, 495, 496, 499, 500, 501, 502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 521, 522, 523, 524, 520, 598, 599".
* `use_jar` - (Optional) Set to true to enable the Cookie Jar. Required for some redirects. Default is false.
* `post_raw` - (Optional) Use to populate the RAW POST data field on the test.
* `final_endpoint` - (Optional) Use to specify the expected Final URL in the testing process.
* `enable_ssl_alert` - (Optional) HTTP Tests only. If enabled, tests will send warnings if the SSL certificate is about to expire. Paid users only. Default is false
* `follow_redirect` - (Optional) Use to specify whether redirects should be followed, set to true to enable. Default is false.
* `dns_server` - (Optional) *Used only for DNS type tests* Hostname or IP of DNS server to use.
* `dns_ip` - (Optional) *Used only for DNS type tests* Comma-separated IPs to compare against the `website_url` resolved value.

## Attributes Reference

The following attribute is exported:

* `test_id` - A unique identifier for the test.

## Import

StatusCake test can be imported using the test id, e.g.

```
tf import statuscake_test.example 123
```