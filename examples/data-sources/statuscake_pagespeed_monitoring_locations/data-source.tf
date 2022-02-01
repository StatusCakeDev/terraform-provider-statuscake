data "statuscake_pagespeed_monitoring_locations" "pagespeed" {
  region_code = "GB"
}

output "pagespeed_monitoring_location_ips" {
  value = toset([for loc in data.statuscake_pagespeed_monitoring_locations.pagespeed.locations : loc.ipv4])
}
