package nexus

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccProviders map[string]terraform.ResourceProvider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"nexus": testAccProvider,
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
	if v := os.Getenv("NEXUS_URL"); v == "" {
		t.Fatalf("NEXUS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("NEXUS_USERNAME"); v == "" {
		t.Fatalf("NEXUS_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("NEXUS_PASSWORD"); v == "" {
		t.Fatalf("NEXUS_PASSWORD must be set for acceptance tests")
	}
}
