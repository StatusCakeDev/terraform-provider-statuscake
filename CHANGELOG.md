## 1.1.0 (Unreleased)
## 1.0.0 (July 02, 2019)

NOTES:

* resource/statuscake_test: `contact_id (int)` has been deprecated, use instead: `contact_group (type: list)` ([#18](https://github.com/terraform-providers/terraform-provider-statuscake/issues/18))
* resource/statuscake_test: `test_tags` has been changed from a CSV string to list of strings ([#18](https://github.com/terraform-providers/terraform-provider-statuscake/issues/18))

FEATURES:

* Add support for Terraform 0.12 ([#32](https://github.com/terraform-providers/terraform-provider-statuscake/issues/32))
* resource/statuscake_test: add `enable_ssl_alert` feature for uptime tests ([#26](https://github.com/terraform-providers/terraform-provider-statuscake/issues/26))
* resource/statuscake_test: add support for Terraform ResourceImporter ([#36](https://github.com/terraform-providers/terraform-provider-statuscake/issues/36))

IMPROVEMENTS:

* resource/statuscake_test: add `contact_group` with multiple contact IDs ([#18](https://github.com/terraform-providers/terraform-provider-statuscake/issues/18), [#34](https://github.com/terraform-providers/terraform-provider-statuscake/issues/34)))

## 0.2.0 (July 27, 2018)

IMPROVEMENTS:

* resource/statuscake_test: Add support for all support library options including basic auth, status codes, custom headers and more. ([#11](https://github.com/terraform-providers/terraform-provider-statuscake/issues/11))


BUG FIXES:

* Fix handling of `contact_id` ([#11](https://github.com/terraform-providers/terraform-provider-statuscake/issues/11))


## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
