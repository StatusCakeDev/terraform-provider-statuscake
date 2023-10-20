resource "statuscake_pagespeed_check" "example_com" {
  check_interval = 300
  name           = "Example"
  region         = "UK"

  alert_config {
    alert_bigger = "5000"
  }

  contact_groups = [
    statuscake_contact_group.operations_team.id
  ]

  monitored_resource {
    address = "https://www.example.com"
  }
}

output "example_com_pagespeed_check_id" {
  value = statuscake_pagespeed_check.example_com.id
}
