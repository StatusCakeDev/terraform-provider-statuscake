resource "statuscake_maintenance_window" "weekends" {
  end             = "2022-01-30T23:59:59Z"
  name            = "Weekends"
  repeat_interval = "1w"
  start           = "2022-01-29T00:00:00Z"
  timezone        = "UTC"

  tags = [
    "production"
  ]

  tests = [
    statuscake_uptime_check.statuscake_com.id,
  ]
}
