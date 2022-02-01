resource "statuscake_contact_group" "operations_team" {
  name     = "Operations Team"
  ping_url = "https://www.example.com"

  email_addresses = [
    "johnsmith@example.com",
    "janesmith@example.com",
  ]
}

output "operations_team_contact_group_id" {
  value = statuscake_contact_group.operations_team.id
}
