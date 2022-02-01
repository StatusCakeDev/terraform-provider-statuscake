resource "statuscake_ssl_check" "example_com" {
  check_interval = 600
  user_agent     = "terraform managed SSL check"

  contact_groups = [
    statuscake_contact_group.operations_team.id
  ]

  alert_config {
    alert_at = [7, 14, 21]

    on_reminder = true
    on_expiry   = true
    on_broken   = false
    on_mixed    = false
  }

  monitored_resource {
    address = "https://www.example.com"
  }
}

output "example_com_ssl_check_id" {
  value = statuscake_ssl_check.example_com.id
}
