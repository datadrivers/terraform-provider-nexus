package security_test

import (
	"fmt"
	"testing"

	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSecurityContentSelector(t *testing.T) {
	var contentSelector security.ContentSelector

	resName := "nexus_security_content_selector.acceptance"
	cs := security.ContentSelector{
		Name:        acctest.RandString(10),
		Description: acctest.RandString(30),
		Expression:  fmt.Sprintf("format == '%s' and path == '%s'", acctest.RandString(15), acctest.RandString(15)),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic content selector
			{
				Config: testAccResourceSecurityContentSelectorConfig(cs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "description", cs.Description),
					resource.TestCheckResourceAttr(resName, "expression", cs.Expression),
					resource.TestCheckResourceAttr(resName, "name", cs.Name),
					testAccCheckContentSelectorResourceExists(resName, &contentSelector),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     contentSelector.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSecurityContentSelectorConfig(cs security.ContentSelector) string {
	return fmt.Sprintf(`
resource "nexus_security_content_selector" "acceptance" {
	description = "%s"
	expression  = "%s"
	name        = "%s"
}
`, cs.Description, cs.Expression, cs.Name)
}

func testAccCheckContentSelectorResourceExists(name string, contentSelector *security.ContentSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client := acceptance.TestAccProvider.Meta().(*nexus.NexusClient)
		result, err := client.Security.ContentSelector.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		*contentSelector = *result

		return nil
	}
}
