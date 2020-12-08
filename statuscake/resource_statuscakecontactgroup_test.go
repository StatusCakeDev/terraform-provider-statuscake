package statuscake

import (
	"fmt"
	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"testing"
)

func TestAccStatusCakeContactGroup_basic(t *testing.T) {
	var contactGroup statuscake.ContactGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccContactGroupCheckDestroy(&contactGroup),
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccContactGroupCheckExists("statuscake_contact_group.exemple", &contactGroup),
					testAccContactGroupCheckAttributes("statuscake_contact_group.exemple", &contactGroup),
				),
			},
		},
	})
}

func TestAccStatusCakeContactGroup_withUpdate(t *testing.T) {
	var contactGroup statuscake.ContactGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccContactGroupCheckDestroy(&contactGroup),
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccContactGroupCheckExists("statuscake_contact_group.exemple", &contactGroup),
					testAccContactGroupCheckAttributes("statuscake_contact_group.exemple", &contactGroup),
				),
			},

			{
				Config: testAccContactGroupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccContactGroupCheckExists("statuscake_contact_group.exemple", &contactGroup),
					testAccContactGroupCheckAttributes("statuscake_contact_group.exemple", &contactGroup),
					resource.TestCheckResourceAttr("statuscake_contact_group.exemple", "group_name", "group"),
					resource.TestCheckResourceAttr("statuscake_contact_group.exemple", "ping_url", "https"),
				),
			},
		},
	})
}

func testAccContactGroupCheckExists(rn string, contactGroup *statuscake.ContactGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ContactGroupID not set")
		}

		client := testAccProvider.Meta().(*statuscake.Client)
		contactGroupId, _ := strconv.Atoi(rs.Primary.ID)

		gotContactGroup, err := statuscake.NewContactGroups(client).Detail(contactGroupId)
		if err != nil {
			return fmt.Errorf("error getting ContactGroup: %s", err)
		}

		*contactGroup = *gotContactGroup

		return nil
	}
}

func testAccContactGroupCheckAttributes(rn string, contactGroup *statuscake.ContactGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[rn].Primary.Attributes

		check := func(key, stateValue, contactGroupValue string) error {
			if contactGroupValue != stateValue {
				return fmt.Errorf("different values for %s in state (%s) and in statuscake (%s)",
					key, stateValue, contactGroupValue)
			}
			return nil
		}

		for key, value := range attrs {
			var err error

			switch key {
			case "contact_id":
				err = check(key, value, strconv.Itoa(contactGroup.ContactID))
			case "desktop_alert":
				err = check(key, value, contactGroup.DesktopAlert)
			case "ping_url":
				err = check(key, value, contactGroup.PingURL)
			case "group_name":
				err = check(key, value, contactGroup.GroupName)
			case "pushover":
				err = check(key, value, contactGroup.Pushover)
			case "boxcar":
				err = check(key, value, contactGroup.Boxcar)
			case "mobiles":
				err = check(key, value, contactGroup.Mobiles)
			case "emails":
				for _, tv := range contactGroup.Emails {
					err = check(key, value, tv)
					if err != nil {
						return err
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

func testAccContactGroupCheckDestroy(contactGroup *statuscake.ContactGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*statuscake.Client)
		_, err := statuscake.NewContactGroups(client).Detail(contactGroup.ContactID)
		if err == nil {
			return fmt.Errorf("contact_group still exists")
		}

		return nil
	}
}

const testAccContactGroupConfig_basic = `
resource "statuscake_contact_group" "exemple" {
	emails= ["aaa","bbb"]
        group_name= "groupname"
        ping_url= "http"
}
`

const testAccContactGroupConfig_update = `
resource "statuscake_contact_group" "exemple" {
         emails= ["aaa","bbb","ccc"]
         group_name= "group"
         ping_url= "https"
}
`
