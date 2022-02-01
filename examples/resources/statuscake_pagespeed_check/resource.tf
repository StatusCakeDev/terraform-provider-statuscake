resource "statuscake_pagespeed_check" "example_com" {
  name           = "Example"
  check_interval = 300
  region         = "UK"

  contact_groups = [
    statuscake_contact_group.operations_team.id
  ]

  alert_config {
    alert_bigger = "5000"
  }

  monitored_resource {
    address = "https://www.example.com"
  }
}

output "example_com_pagespeed_check_id" {
  value = statuscake_pagespeed_check.example_com.id
}
