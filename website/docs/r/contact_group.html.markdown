---
layout: "statuscake"
page_title: "StatusCake: statuscake_contact_group"
sidebar_current: "docs-statuscake-contact_group"
description: |-
  The statuscake_contact_group resource allows StatusCake contact groups to be managed by Terraform.
---

# statuscake\_contact_group

The contact_group resource allows StatusCake contact groups to be managed by Terraform.

## Example Usage

```hcl
resource "statuscake_contact_group" "exemple" {
	emails= ["email1","email2"]
        group_name= "group name"
        ping_url= "url"
}
```

## Argument Reference

The following arguments are supported:

* `desktop_alert` - (Required) Set to 1 To Enable Desktop Alerts
* `ping_url` - (Optional) 
* `group_name` - (Optional) The internal Group Name
* `pushover` - (Optional) A Pushover Account Key
* `boxcar` - (Optional) A Boxcar API Key
* `mobiles` - (Optional) Comma Seperated List of International Format Cell Numbers
* `emails` - (Optional) List of Emails To Alert.

## Attributes Reference

The following attribute is exported:

* `contact_id` - A unique identifier for the contact group.

## Import

StatusCake contact groups can be imported using the contact group id, e.g.

```
tf import statuscake_contact_group.example 123
```
