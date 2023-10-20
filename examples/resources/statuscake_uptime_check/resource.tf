resource "statuscake_uptime_check" "example_com" {
  check_interval = 30
  confirmation   = 3
  name           = "Example"
  trigger_rate   = 10

  contact_groups = [
    statuscake_contact_group.operations_team.id
  ]

  http_check {
    enable_cookies   = true
    follow_redirects = true
    timeout          = 20
    user_agent       = "terraform managed uptime check"
    validate_ssl     = true

    basic_authentication {
      password = "password"
      username = "username"
    }

    content_matchers {
      content         = "Welcome"
      include_headers = true
      matcher         = "CONTAINS_STRING"
    }

    request_headers = {
      Authorization = "bearer 123456"
    }

    status_codes = [
      "202",
      "404",
      "405",
    ]
  }

  monitored_resource {
    address = "https://www.example.com"
  }

  regions = [
    "london",
    "london",
    "paris",
  ]

  tags = [
    "production",
  ]
}

output "example_com_uptime_check_id" {
  value = statuscake_uptime_check.example_com.id
}
