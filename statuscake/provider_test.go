package statuscake

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"statuscake": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("STATUSCAKE_USERNAME"); v == "" {
		t.Fatal("STATUSCAKE_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("STATUSCAKE_APIKEY"); v == "" {
		t.Fatal("STATUSCAKE_APIKEY must be set for acceptance tests")
	}
	if v := os.Getenv("STATUSCAKE_TEST_CONTACT_GROUP_ID"); v == "" {
		t.Fatal("STATUSCAKE_TEST_CONTACT_GROUP_ID must be set for acceptance tests")
	}
}
