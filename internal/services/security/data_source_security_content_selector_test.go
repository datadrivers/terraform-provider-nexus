package security_test

import (
	"fmt"
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityContentSelector(t *testing.T) {
	dataSourceName := "data.nexus_security_content_selector.acceptance"

	cs := security.ContentSelector{
		Name:        acctest.RandString(10),
		Description: acctest.RandString(30),
		Expression:  fmt.Sprintf("format == '%s' and path == '%s'", acctest.RandString(15), acctest.RandString(15)),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityContentSelectorConfig(cs) + testAccDataSourceSecurityContentSelectorConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", cs.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", cs.Description),
					resource.TestCheckResourceAttr(dataSourceName, "expression", cs.Expression),
				),
			},
		},
	})
}

func testAccDataSourceSecurityContentSelectorConfig() string {
	return `
data "nexus_security_content_selector" "acceptance" {
	name = nexus_security_content_selector.acceptance.name
}
`
}
