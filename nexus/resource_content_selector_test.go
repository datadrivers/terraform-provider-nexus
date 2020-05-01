package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccContentSelector(t *testing.T) {
	t.Parallel()

	var contentSelector nexus.ContentSelector

	contentSelectorName := acctest.RandString(10)
	contentSelectorDescription := acctest.RandString(30)
	contentSelectorExpression := fmt.Sprintf("format == \\\"%s\\\" and path == \\\"%s\\\"", acctest.RandString(15), acctest.RandString(15))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic content selector
			{
				Config: testAccContentSelectorResource(contentSelectorName, contentSelectorDescription, contentSelectorExpression),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_content_selector.acceptance", "description", contentSelectorDescription),
					//resource.TestCheckResourceAttr("nexus_content_selector.acceptance", "expression", contentSelectorExpression),
					resource.TestCheckResourceAttr("nexus_content_selector.acceptance", "name", contentSelectorName),
					testAccCheckContentSelectorResourceExists("nexus_content_selector.acceptance", &contentSelector),
				),
			},
			{
				ResourceName:      "nexus_content_selector.acceptance",
				ImportStateId:     contentSelectorName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContentSelectorResource(name string, description string, expression string) string {

	return fmt.Sprintf(`
	resource "nexus_content_selector" "acceptance" {
		name   = "%s"
		description = "%s"
		expression = "%s"
	}
	`, name, description, expression)
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
