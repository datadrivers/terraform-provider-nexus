package other_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/capability"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceCapability() capability.CapabilityCreate {
	return capability.CapabilityCreate{
		Type:    "webhook.global",
		Notes:   fmt.Sprintf("TERRAFORM_TEST_%s", acctest.RandString(10)),
		Enabled: false, // Keep disabled to avoid triggering actual webhooks
		Properties: map[string]string{
			"url":    "https://example.com/webhook",
			"secret": "test-secret-" + acctest.RandString(10),
			"names":  "audit", // Required: valid event types are audit, repository, etc.
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
					resource.TestCheckResourceAttr(resName, "properties.url", cap.Properties["url"]),
					resource.TestCheckResourceAttr(resName, "properties.secret", cap.Properties["secret"]),
					resource.TestCheckResourceAttr(resName, "properties.names", cap.Properties["names"]),
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
	// Note: ID must be included in the PUT body for updates to work

	resName := "nexus_capability.acceptance"

	cap := testAccResourceCapability()
	capUpdated := capability.CapabilityCreate{
		Type:    cap.Type,
		Notes:   fmt.Sprintf("updated-capability-%s", acctest.RandString(10)),
		Enabled: true, // Changed from false to true to test enable toggle
		Properties: map[string]string{
			"url":    "https://updated.example.com/webhook",
			"secret": "updated-secret-" + acctest.RandString(10),
			"names":  "repository", // Different event type for update test
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
					resource.TestCheckResourceAttr(resName, "properties.url", cap.Properties["url"]),
					resource.TestCheckResourceAttr(resName, "properties.secret", cap.Properties["secret"]),
					resource.TestCheckResourceAttr(resName, "properties.names", cap.Properties["names"]),
				),
			},
			{
				Config: testAccResourceCapabilityConfig(capUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "type", capUpdated.Type),
					resource.TestCheckResourceAttr(resName, "notes", capUpdated.Notes),
					resource.TestCheckResourceAttr(resName, "enabled", fmt.Sprintf("%t", capUpdated.Enabled)),
					resource.TestCheckResourceAttr(resName, "properties.url", capUpdated.Properties["url"]),
					resource.TestCheckResourceAttr(resName, "properties.secret", capUpdated.Properties["secret"]),
					resource.TestCheckResourceAttr(resName, "properties.names", capUpdated.Properties["names"]),
				),
			},
		},
	})
}

func TestAccResourceCapabilityWithProperties(t *testing.T) {
	resName := "nexus_capability.acceptance"

	cap := capability.CapabilityCreate{
		Type:    "webhook.global",
		Notes:   fmt.Sprintf("TERRAFORM_TEST_props_%s", acctest.RandString(10)),
		Enabled: false,
		Properties: map[string]string{
			"url":    "https://webhook-test.example.com/hook",
			"secret": "test-props-secret-" + acctest.RandString(10),
			"names":  "audit", // Valid event type
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
					resource.TestCheckResourceAttr(resName, "properties.url", cap.Properties["url"]),
					resource.TestCheckResourceAttr(resName, "properties.secret", cap.Properties["secret"]),
					resource.TestCheckResourceAttr(resName, "properties.names", cap.Properties["names"]),
				),
			},
		},
	})
}

func testAccResourceCapabilityConfig(cap capability.CapabilityCreate) string {
	propertiesHCL := ""
	if len(cap.Properties) > 0 {
		propertiesHCL = "\n\tproperties = {\n"
		for k, v := range cap.Properties {
			propertiesHCL += fmt.Sprintf("\t\t%s = \"%s\"\n", k, v)
		}
		propertiesHCL += "\t}"
	}

	return fmt.Sprintf(`
resource "nexus_capability" "acceptance" {
	type    = "%s"
	notes   = "%s"
	enabled = %t%s
}
`, cap.Type, cap.Notes, cap.Enabled, propertiesHCL)
}

func testAccResourceCapabilityConfigWithMultipleProperties(cap capability.CapabilityCreate) string {
	propertiesHCL := ""
	if len(cap.Properties) > 0 {
		propertiesHCL = "\n\tproperties = {\n"
		for k, v := range cap.Properties {
			propertiesHCL += fmt.Sprintf("\t\t%s = \"%s\"\n", k, v)
		}
		propertiesHCL += "\t}"
	}

	return fmt.Sprintf(`
resource "nexus_capability" "acceptance" {
	type    = "%s"
	notes   = "%s"
	enabled = %t%s
}
`, cap.Type, cap.Notes, cap.Enabled, propertiesHCL)
}
