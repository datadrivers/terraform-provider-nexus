package security_test

import (
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePrivilegeWildcard(t *testing.T) {
	dataSourceName := "data.nexus_privilege_wildcard.acceptance"

	privilege := security.PrivilegeWildcard{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Pattern:     "nexus:*",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeWildcardConfig(privilege) + testAccDataSourcePrivilegeWildcardConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", privilege.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", privilege.Description),
					resource.TestCheckResourceAttr(dataSourceName, "pattern", privilege.Pattern),
				),
			},
		},
	})
}

func testAccDataSourcePrivilegeWildcardConfig() string {
	return `
	data "nexus_privilege_wildcard" "acceptance" {
		name = nexus_privilege_wildcard.acceptance.name
	}`
}
