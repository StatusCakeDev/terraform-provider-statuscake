data "statuscake_uptime_monitoring_locations" "uptime" {
  region_code = "GBR"
}

output "uptime_monitoring_location_ips" {
  value = toset([for loc in data.statuscake_uptime_monitoring_locations.uptime.locations : loc.ipv4])
}
