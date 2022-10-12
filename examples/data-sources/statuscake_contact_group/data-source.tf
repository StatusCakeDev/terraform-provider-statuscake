data "statuscake_contact_group" "developers" {
  id = "123456"
}

output "developers_contact_group_name" {
  value = data.statuscake_contact_group.developers.name
}
