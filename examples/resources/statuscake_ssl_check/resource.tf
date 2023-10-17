resource "statuscake_ssl_check" "example_com" {
  check_interval = 600
  user_agent     = "terraform managed SSL check"

  alert_config {
    alert_at    = [7, 14, 21]
    on_broken   = false
    on_expiry   = true
    on_mixed    = false
    on_reminder = true
  }

  contact_groups = [
    statuscake_contact_group.operations_team.id
  ]

  monitored_resource {
    address = "https://www.example.com"
  }
}

output "example_com_ssl_check_id" {
  value = statuscake_ssl_check.example_com.id
}
