package other_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/capability"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceCapability() capability.CapabilityCreate {
	return capability.CapabilityCreate{
		Type:    "OutreachManagementCapability",
		Notes:   fmt.Sprintf("TERRAFORM_TEST_%s", acctest.RandString(10)),
		Enabled: false, // Disabled to avoid side effects
		Properties: map[string]string{
			"baseUrl":      "https://links.sonatype.com/products/nexus/outreach",
			"alwaysRemote": "false",
		},
	}
}

func TestAccResourceCapability(t *testing.T) {
	resName := "nexus_capability.acceptance"

	cap := testAccResourceCapability()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCapabilityConfig(cap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "type", cap.Type),
					resource.TestCheckResourceAttr(resName, "notes", cap.Notes),
					resource.TestCheckResourceAttr(resName, "enabled", fmt.Sprintf("%t", cap.Enabled)),
					resource.TestCheckResourceAttr(resName, "properties.baseUrl", "https://links.sonatype.com/products/nexus/outreach"),
					resource.TestCheckResourceAttr(resName, "properties.alwaysRemote", "false"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceCapabilityUpdate(t *testing.T) {
	// SKIP: OutreachManagementCapability has known Nexus bugs with updates
	// causing NullPointerException. This is a Nexus server issue.
	// See: java.lang.NullPointerException: Cannot invoke "...CapabilityReference.context()"
	t.Skip("Update test skipped due to known Nexus bug with OutreachManagementCapability updates")

	resName := "nexus_capability.acceptance"

	cap := testAccResourceCapability()
	capUpdated := capability.CapabilityCreate{
		Type:    cap.Type,
		Notes:   fmt.Sprintf("updated-capability-%s", acctest.RandString(10)),
		Enabled: false, // Keep same enabled state
		Properties: map[string]string{
			"baseUrl":      "https://links.sonatype.com/products/nexus/outreach",
			"alwaysRemote": "false",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCapabilityConfig(cap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "type", cap.Type),
					resource.TestCheckResourceAttr(resName, "notes", cap.Notes),
					resource.TestCheckResourceAttr(resName, "enabled", fmt.Sprintf("%t", cap.Enabled)),
				),
			},
			{
				Config: testAccResourceCapabilityConfig(capUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "type", capUpdated.Type),
					resource.TestCheckResourceAttr(resName, "notes", capUpdated.Notes),
					resource.TestCheckResourceAttr(resName, "enabled", fmt.Sprintf("%t", capUpdated.Enabled)),
				),
			},
		},
	})
}

func TestAccResourceCapabilityWithProperties(t *testing.T) {
	resName := "nexus_capability.acceptance"

	cap := capability.CapabilityCreate{
		Type:    "OutreachManagementCapability",
		Notes:   fmt.Sprintf("TERRAFORM_TEST_props_%s", acctest.RandString(10)),
		Enabled: false,
		Properties: map[string]string{
			"baseUrl":      "https://links.sonatype.com/products/nexus/outreach",
			"alwaysRemote": "false",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCapabilityConfigWithMultipleProperties(cap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "type", cap.Type),
					resource.TestCheckResourceAttr(resName, "notes", cap.Notes),
					resource.TestCheckResourceAttr(resName, "enabled", fmt.Sprintf("%t", cap.Enabled)),
					resource.TestCheckResourceAttr(resName, "properties.baseUrl", "https://links.sonatype.com/products/nexus/outreach"),
					resource.TestCheckResourceAttr(resName, "properties.alwaysRemote", "false"),
				),
			},
		},
	})
}

func testAccResourceCapabilityConfig(cap capability.CapabilityCreate) string {
	return fmt.Sprintf(`
resource "nexus_capability" "acceptance" {
	type    = "%s"
	notes   = "%s"
	enabled = %t
	properties = {
		baseUrl      = "%s"
		alwaysRemote = "%s"
	}
}
`, cap.Type, cap.Notes, cap.Enabled, cap.Properties["baseUrl"], cap.Properties["alwaysRemote"])
}

func testAccResourceCapabilityConfigWithMultipleProperties(cap capability.CapabilityCreate) string {
	return fmt.Sprintf(`
resource "nexus_capability" "acceptance" {
	type    = "%s"
	notes   = "%s"
	enabled = %t
	properties = {
		baseUrl      = "%s"
		alwaysRemote = "%s"
	}
}
`, cap.Type, cap.Notes, cap.Enabled, cap.Properties["baseUrl"], cap.Properties["alwaysRemote"])
}
