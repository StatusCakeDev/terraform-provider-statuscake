## 1.0.0 (Unreleased)

NOTES:

* resource/statuscake_test: `contact_id (int)` has been deprecated, use instead: `contact_group (type: list)`
* resource:statuscake_test: `test_tags` has been changed from a CSV string to list of strings


IMPROVEMENTS:

* resource/statuscake_test: add `contact_group` with multiple contact IDs [GH-8]


## 0.2.0 (July 27, 2018)

IMPROVEMENTS:

* resource/statuscake_test: Add support for all support library options including basic auth, status codes, custom headers and more. ([#11](https://github.com/terraform-providers/terraform-provider-statuscake/issues/11))


BUG FIXES:

* Fix handling of `contact_id` ([#11](https://github.com/terraform-providers/terraform-provider-statuscake/issues/11))


## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
