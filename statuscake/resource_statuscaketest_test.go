package statuscake

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccStatusCake_basic(t *testing.T) {
	var test statuscake.Test

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTestCheckDestroy(&test),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplate(testAccTestConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("statuscake_test.google", &test),
					testAccTestCheckAttributes("statuscake_test.google", &test),
				),
			},
		},
	})
}

func TestAccStatusCake_basic_deprecated_contact_ID(t *testing.T) {
	var test statuscake.Test

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTestCheckDestroy(&test),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplate(testAccTestConfig_deprecated),
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("statuscake_test.google", &test),
					testAccTestCheckAttributes("statuscake_test.google", &test),
				),
			},
		},
	})
}

func TestAccStatusCake_tcp(t *testing.T) {
	var test statuscake.Test

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTestCheckDestroy(&test),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplate(testAccTestConfig_tcp),
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("statuscake_test.google", &test),
					testAccTestCheckAttributes("statuscake_test.google", &test),
				),
			},
		},
	})
}

func TestAccStatusCake_dns(t *testing.T) {
	var test statuscake.Test

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTestCheckDestroy(&test),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplate(testAccTestConfig_dns),
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("statuscake_test.google", &test),
					testAccTestCheckAttributes("statuscake_test.google", &test),
				),
			},
		},
	})
}

func TestAccStatusCake_withUpdate(t *testing.T) {
	var test statuscake.Test

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTestCheckDestroy(&test),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplate(testAccTestConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("statuscake_test.google", &test),
				),
			},

			{
				Config: testAccTestConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("statuscake_test.google", &test),
					testAccTestCheckAttributes("statuscake_test.google", &test),
					resource.TestCheckResourceAttr("statuscake_test.google", "check_rate", "500"),
					resource.TestCheckResourceAttr("statuscake_test.google", "paused", "true"),
					resource.TestCheckResourceAttr("statuscake_test.google", "timeout", "40"),
					resource.TestCheckResourceAttr("statuscake_test.google", "confirmations", "0"),
					resource.TestCheckResourceAttr("statuscake_test.google", "trigger_rate", "20"),
					resource.TestCheckResourceAttr("statuscake_test.google", "custom_header", "{ \"Content-Type\": \"application/x-www-form-urlencoded\" }"),
					resource.TestCheckResourceAttr("statuscake_test.google", "user_agent", "string9988"),
					resource.TestCheckResourceAttr("statuscake_test.google", "status", "Up"),
					resource.TestCheckResourceAttr("statuscake_test.google", "uptime", "0"),
					resource.TestCheckResourceAttr("statuscake_test.google", "node_locations.#", "3"),
					resource.TestCheckResourceAttr("statuscake_test.google", "ping_url", "string8410"),
					resource.TestCheckResourceAttr("statuscake_test.google", "basic_user", "string27052"),
					resource.TestCheckResourceAttr("statuscake_test.google", "basic_pass", "string5659"),
					resource.TestCheckResourceAttr("statuscake_test.google", "public", "0"),
					resource.TestCheckResourceAttr("statuscake_test.google", "logo_image", "string21087"),
					resource.TestCheckResourceAttr("statuscake_test.google", "branding", "25875"),
					resource.TestCheckResourceAttr("statuscake_test.google", "website_host", "string32368"),
					resource.TestCheckResourceAttr("statuscake_test.google", "virus", "1"),
					resource.TestCheckResourceAttr("statuscake_test.google", "find_string", "string15212"),
					resource.TestCheckResourceAttr("statuscake_test.google", "do_not_find", "false"),
					resource.TestCheckResourceAttr("statuscake_test.google", "real_browser", "1"),
					resource.TestCheckResourceAttr("statuscake_test.google", "test_tags.#", "1"),
					resource.TestCheckResourceAttr("statuscake_test.google", "status_codes", "string23065"),
					resource.TestCheckResourceAttr("statuscake_test.google", "use_jar", "1"),
					resource.TestCheckResourceAttr("statuscake_test.google", "post_raw", "string32096"),
					resource.TestCheckResourceAttr("statuscake_test.google", "final_endpoint", "string10781"),
					resource.TestCheckResourceAttr("statuscake_test.google", "enable_ssl_alert", "false"),
					resource.TestCheckResourceAttr("statuscake_test.google", "follow_redirect", "true"),
				),
			},
		},
	})
}

func testAccTestCheckExists(rn string, test *statuscake.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TestID not set")
		}

		client := testAccProvider.Meta().(*statuscake.Client)
		testId, parseErr := strconv.Atoi(rs.Primary.ID)
		if parseErr != nil {
			return fmt.Errorf("error in statuscake test CheckExists: %s", parseErr)
		}

		gotTest, err := client.Tests().Detail(testId)
		if err != nil {
			return fmt.Errorf("error getting test: %s", err)
		}

		*test = *gotTest

		return nil
	}
}

func testAccTestCheckAttributes(rn string, test *statuscake.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[rn].Primary.Attributes

		check := func(key, stateValue, testValue string) error {
			if testValue != stateValue {
				return fmt.Errorf("different values for %s in state (%s) and in statuscake (%s)",
					key, stateValue, testValue)
			}
			return nil
		}

		for key, value := range attrs {
			var err error

			switch key {
			case "website_name":
				err = check(key, value, test.WebsiteName)
			case "website_url":
				err = check(key, value, test.WebsiteURL)
			case "check_rate":
				err = check(key, value, strconv.Itoa(test.CheckRate))
			case "test_type":
				err = check(key, value, test.TestType)
			case "paused":
				err = check(key, value, strconv.FormatBool(test.Paused))
			case "timeout":
				err = check(key, value, strconv.Itoa(test.Timeout))
			case "contact_group":
				for _, tv := range test.ContactGroup {
					err = check(key, value, tv)
					if err != nil {
						return err
					}
				}
			case "confirmations":
				err = check(key, value, strconv.Itoa(test.Confirmation))
			case "trigger_rate":
				err = check(key, value, strconv.Itoa(test.TriggerRate))
			case "custom_header":
				err = check(key, value, test.CustomHeader)
			case "node_locations":
				for _, tv := range test.NodeLocations {
					err = check(key, value, tv)
					if err != nil {
						return err
					}
				}
			case "public":
				err = check(key, value, strconv.Itoa(test.Public))
			case "logo_image":
				err = check(key, value, test.LogoImage)
			case "find_string":
				err = check(key, value, test.FindString)
			case "do_not_find":
				err = check(key, value, strconv.FormatBool(test.DoNotFind))
			case "status_codes":
				err = check(key, value, test.StatusCodes)
			case "use_jar":
				err = check(key, value, strconv.Itoa(test.UseJar))
			case "post_raw":
				err = check(key, value, test.PostRaw)
			case "final_endpoint":
				err = check(key, value, test.FinalEndpoint)
			case "enable_ssl_alert":
				err = check(key, value, strconv.FormatBool(test.EnableSSLAlert))
			case "follow_redirect":
				err = check(key, value, strconv.FormatBool(test.FollowRedirect))
			case "dns_server":
				err = check(key, value, test.DNSServer)
			case "dns_ip":
				err = check(key, value, test.DNSIP)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccTestCheckDestroy(test *statuscake.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*statuscake.Client)
		err := client.Tests().Delete(test.TestID)
		if err == nil {
			return fmt.Errorf("test still exists")
		}

		return nil
	}
}

func interpolateTerraformTemplate(template string) string {
	testContactGroupId := "43402"

	if v := os.Getenv("STATUSCAKE_TEST_CONTACT_GROUP_ID"); v != "" {
		testContactGroupId = v
	}

	return fmt.Sprintf(template, testContactGroupId)
}

const testAccTestConfig_basic = `
resource "statuscake_test" "google" {
	website_name = "google.com"
	website_url = "www.google.com"
	test_type = "HTTP"
	check_rate = 300
	timeout = 10
	contact_group = ["%s"]
	confirmations = 1
	trigger_rate = 10
}
`
const testAccTestConfig_deprecated = `
resource "statuscake_test" "google" {
	website_name = "google.com"
	website_url = "www.google.com"
	test_type = "HTTP"
	check_rate = 300
	timeout = 10
	contact_id = %s
	confirmations = 1
	trigger_rate = 10
}
`
const testAccTestConfig_update = `
resource "statuscake_test" "google" {
	website_name = "google.com"
	website_url = "www.google.com"
	test_type = "HTTP"
	check_rate = 500
	paused = true
	trigger_rate = 20
	custom_header = "{ \"Content-Type\": \"application/x-www-form-urlencoded\" }"
	user_agent = "string9988"
	node_locations = [ "string16045", "string19741", "string12122" ]
	ping_url = "string8410"
	basic_user = "string27052"
	basic_pass = "string5659"
	public = 0
	logo_image = "string21087"
	branding = 25875
	website_host = "string32368"
	virus = 1
	find_string = "string15212"
	do_not_find = false
	real_browser = 1
	test_tags = ["string8191"]
	status_codes = "string23065"
	use_jar = 1
	post_raw = "string32096"
	final_endpoint = "string10781"
	enable_ssl_alert = false
	follow_redirect = true
}
`

const testAccTestConfig_tcp = `
resource "statuscake_test" "google" {
	website_name = "google.com"
	website_url = "www.google.com"
	test_type = "TCP"
	check_rate = 300
	timeout = 10
	contact_group = ["%s"]
	confirmations = 1
	port = 80
}
`
const testAccTestConfig_dns = `
resource "statuscake_test" "google" {
	website_url = "dns.google"
	dns_server = "1.1.1.1"
	dns_ip = "8.8.8.8,8.8.4.4"
	test_type = "DNS"
	check_rate = 300
	timeout = 10
	contact_group = ["%s"]
	confirmations = 1
}
`
