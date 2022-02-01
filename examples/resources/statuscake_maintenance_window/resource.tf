resource "statuscake_maintenance_window" "weekends" {
  name     = "Weekends"
  timezone = "UTC"

  start           = "2022-01-29T00:00:00Z"
  end             = "2022-01-30T23:59:59Z"
  repeat_interval = "1w"

  tags = [
    "production"
  ]

  tests = [
    statuscake_uptime_check.statuscake_com.id,
  ]
}
