package security_test

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePrivilegeApplication(t *testing.T) {
	dataSourceName := "data.nexus_privilege_application.acceptance"

	privilege := security.PrivilegeApplication{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Domain:      acctest.RandString(20),
		Actions:     []security.SecurityPrivilegeApplicationActions{"DELETE"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeApplicationConfig(privilege) + testAccDataSourcePrivilegeApplicationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", privilege.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", privilege.Description),
					resource.TestCheckResourceAttr(dataSourceName, "domain", privilege.Domain),
				),
			},
		},
	})
}

func testAccDataSourcePrivilegeApplicationConfig() string {
	return `
	data "nexus_privilege_application" "acceptance" {
		name = nexus_privilege_application.acceptance.name
	}`
}
