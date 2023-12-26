---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "statuscake_heartbeat_check Resource - terraform-provider-statuscake"
subcategory: ""
description: |-
  
---

# statuscake_heartbeat_check (Resource)



## Example Usage

```terraform
resource "statuscake_uptime_check" "example" {
  name   = "Example"
  period = 1800

  contact_groups = [
    statuscake_contact_group.operations_team.id
  ]

  tags = [
    "production",
  ]
}

output "example_heartbeat_check_id" {
  value = statuscake_heartbeat_check.example.id
}

output "example_heartbeat_check_url" {
  value = statuscake_heartbeat_check.example.check_url
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the check
- `period` (Number) Number of seconds since the last ping before the check is considered down.

### Optional

- `contact_groups` (Set of String) List of contact group IDs
- `monitored_resource` (Block List, Max: 1) Monitored resource configuration block. This describes the server under test (see [below for nested schema](#nestedblock--monitored_resource))
- `paused` (Boolean) Whether the check should be run
- `tags` (Set of String) List of tags

### Read-Only

- `check_url` (String) URL of the heartbeat check
- `id` (String) The ID of this resource.

<a id="nestedblock--monitored_resource"></a>
### Nested Schema for `monitored_resource`

Optional:

- `host` (String) Name of the hosting provider

## Import

Import is supported using the following syntax:

```shell
terraform import statuscake_heartbeat_check.example_com 1234
```