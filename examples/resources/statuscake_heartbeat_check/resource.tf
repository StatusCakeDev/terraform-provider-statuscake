resource "statuscake_heartbeat_check" "example" {
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
