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

func TestAccStatusCakeSsl_basic(t *testing.T) {
	var ssl statuscake.Ssl

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSslCheckDestroy(&ssl),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateSsl(testAccSslConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccSslCheckExists("statuscake_ssl.exemple", &ssl),
					testAccSslCheckAttributes("statuscake_ssl.exemple", &ssl),
				),
			},
		},
	})
}

func TestAccStatusCakeSsl_withUpdate(t *testing.T) {
	var ssl statuscake.Ssl

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSslCheckDestroy(&ssl),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateSsl(testAccSslConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccSslCheckExists("statuscake_ssl.exemple", &ssl),
					testAccSslCheckAttributes("statuscake_ssl.exemple", &ssl),
				),
			},

			{
				Config: testAccSslConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccSslCheckExists("statuscake_ssl.exemple", &ssl),
					testAccSslCheckAttributes("statuscake_ssl.exemple", &ssl),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "check_rate", "86400"),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "domain", "https://www.exemple.com"),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "contact_group_c", ""),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "alert_at", "18,8,2019"),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "alert_reminder", "false"),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "alert_expiry", "false"),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "alert_broken", "true"),
					resource.TestCheckResourceAttr("statuscake_ssl.exemple", "alert_mixed", "false"),
				),
			},
		},
	})
}

func testAccSslCheckExists(rn string, ssl *statuscake.Ssl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SslID not set")
		}

		client := testAccProvider.Meta().(*statuscake.Client)
		sslId := rs.Primary.ID

		gotSsl, err := statuscake.NewSsls(client).Detail(sslId)
		if err != nil {
			return fmt.Errorf("error getting ssl: %s", err)
		}

		*ssl = *gotSsl

		return nil
	}
}

func testAccSslCheckAttributes(rn string, ssl *statuscake.Ssl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[rn].Primary.Attributes

		check := func(key, stateValue, sslValue string) error {
			if sslValue != stateValue {
				return fmt.Errorf("different values for %s in state (%s) and in statuscake (%s)",
					key, stateValue, sslValue)
			}
			return nil
		}

		for key, value := range attrs {
			var err error

			switch key {
			case "domain":
				err = check(key, value, ssl.Domain)
			case "contact_groups_c":
				err = check(key, value, ssl.ContactGroupsC)
			case "checkrate":
				err = check(key, value, strconv.Itoa(ssl.Checkrate))
			case "alert_at":
				err = check(key, value, ssl.AlertAt)
			case "alert_reminder":
				err = check(key, value, strconv.FormatBool(ssl.AlertReminder))
			case "alert_expiry":
				err = check(key, value, strconv.FormatBool(ssl.AlertExpiry))
			case "alert_broken":
				err = check(key, value, strconv.FormatBool(ssl.AlertBroken))
			case "alert_mixed":
				err = check(key, value, strconv.FormatBool(ssl.AlertMixed))
			case "paused":
				err = check(key, value, strconv.FormatBool(ssl.Paused))
			case "issuer_cn":
				err = check(key, value, ssl.IssuerCn)
			case "contact_groups":
				for _, tv := range ssl.ContactGroups {
					err = check(key, value, tv)
					if err != nil {
						return err
					}
				}
			case "cert_score":
				err = check(key, value, ssl.CertScore)
			case "cert_status":
				err = check(key, value, ssl.CertStatus)
			case "cipher":
				err = check(key, value, ssl.Cipher)
			case "valid_from_utc":
				err = check(key, value, ssl.ValidFromUtc)
			case "valid_until_utc":
				err = check(key, value, ssl.ValidUntilUtc)
			case "last_reminder":
				err = check(key, value, strconv.Itoa(ssl.LastReminder))
			case "last_updated_utc":
				err = check(key, value, ssl.LastUpdatedUtc)
			case "flags":
				for _, tv := range ssl.Flags {
					err = check(key, value, strconv.FormatBool(tv))
					if err != nil {
						return err
					}
				}

			case "mixed_content":
				for _, tv := range ssl.MixedContent {
					for _, tv2 := range tv {
						err = check(key, value, tv2)
						if err != nil {
							return err
						}
					}
				}
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccSslCheckDestroy(ssl *statuscake.Ssl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*statuscake.Client)
		_, err := statuscake.NewSsls(client).Detail(ssl.ID)
		if err == nil {
			return fmt.Errorf("ssl still exists")
		}

		return nil
	}
}

func interpolateTerraformTemplateSsl(template string) string {
	sslContactGroupId := "43402"

	if v := os.Getenv("STATUSCAKE_SSL_CONTACT_GROUP_ID"); v != "" {
		sslContactGroupId = v
	}

	return fmt.Sprintf(template, sslContactGroupId)
}

const testAccSslConfig_basic = `
resource "statuscake_ssl" "exemple" {
	domain = "https://www.exemple.com"
	contact_groups_c = "%s"
        checkrate = 3600
        alert_at = "18,7,2019"
        alert_reminder = true
	alert_expiry = true
        alert_broken = false
        alert_mixed = true
}
`

const testAccSslConfig_update = `
resource "statuscake_ssl" "exemple" {
	domain = "https://www.exemple.com"
        contact_groups_c = ""
        checkrate = 86400 
        alert_at = "18,8,2019"
        alert_reminder = false
	alert_expiry = false
        alert_broken = true
        alert_mixed = false
}
`
