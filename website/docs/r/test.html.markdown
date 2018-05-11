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
  website_name = "google.com"
  website_url  = "www.google.com"
  test_type    = "HTTP"
  check_rate   = 300
  contact_id   = 12345
}
```

## Argument Reference

The following arguments are supported:

* `website_name` - (Required) This is the name of the test and the website to be monitored.
* `website_url` - (Required) The URL of the website to be monitored
* `check_rate` - (Optional) Test check rate in seconds. Defaults to 300
* `contact_id` - (Optional) The id of the contact group to be add to the test.  Each test can have only one.
* `test_type` - (Required) The type of Test. Either HTTP or TCP
* `paused` - (Optional) Whether or not the test is paused. Defaults to false.
* `timeout` - (Optional) The timeout of the test in seconds.
* `confirmations` - (Optional) The number of confirmation servers to use in order to detect downtime. Defaults to 0.
* `port` - (Optional) The port to use when specifying a TCP test.
* `trigger_rate` - (Optional) The number of minutes to wait before sending an alert. Default is `5`.
* `custom_header` - (Optional) Custom HTTP header, must be supplied as JSON
* `user_agent` - (Optional) Test with a custom user agent set.
* `use_jar` - (Optional) Set to true to enable the Cookie Jar. Required for some redirects. Default is false.
* `post_raw` - (Optional) Use to populate the RAW POST data field on the test.
* `find_string` - (Optional) A string that should either be found or not found.
* `final_endpoint` - (Optional) Use to specify the expected Final URL in the testing process
* `follow_redirect` - (Optional) Use to specify whether redirects should be followed, set to true to enable. Default is false.
* `status_codes` - (Optional) Comma Seperated List of StatusCodes to Trigger Error on. Defaults are "204, 205, 206, 303, 400, 401, 403, 404, 405, 406, 408, 410, 413, 444, 429, 494, 495, 496, 499, 500, 501, 502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 521, 522, 523, 524, 520, 598, 599"


## Attributes Reference

The following attribute is exported:

* `test_id` - A unique identifier for the test.
