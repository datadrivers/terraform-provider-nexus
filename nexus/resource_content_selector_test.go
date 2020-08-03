package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceContentSelector(t *testing.T) {
	var contentSelector nexus.ContentSelector

	resName := "nexus_content_selector.acceptance"
	cs := nexus.ContentSelector{
		Name:        acctest.RandString(10),
		Description: acctest.RandString(30),
		Expression:  fmt.Sprintf("format == \\\"%s\\\" and path == \\\"%s\\\"", acctest.RandString(15), acctest.RandString(15)),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic content selector
			{
				Config: testAccResourceContentSelectorConfig(cs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "description", cs.Description),
					//resource.TestCheckResourceAttr("nexus_content_selector.acceptance", "expression", contentSelectorExpression),
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

func testAccResourceContentSelectorConfig(cs nexus.ContentSelector) string {
	return fmt.Sprintf(`
resource "nexus_content_selector" "acceptance" {
	description = "%s"
	expression  = "%s"
	name        = "%s"
}
`, cs.Description, cs.Expression, cs.Name)
}

func testAccCheckContentSelectorResourceExists(name string, contentSelector *nexus.ContentSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client := testAccProvider.Meta().(nexus.Client)
		result, err := client.ContentSelectorRead(rs.Primary.ID)
		if err != nil {
			return err
		}

		*contentSelector = *result

		return nil
	}
}
