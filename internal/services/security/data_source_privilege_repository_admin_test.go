package security_test

import (
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePrivilegeRepositoryAdmin(t *testing.T) {
	dataSourceName := "data.nexus_privilege_repository_admin.acceptance"

	privilege := security.PrivilegeRepositoryAdmin{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Actions:     []security.SecurityPrivilegeRepositoryAdminActions{"ADD", "READ", "DELETE", "BROWSE", "EDIT"},
		Repository:  acctest.RandString(20),
		Format:      "helm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeRepositoryAdminConfig(privilege) + testAccDataSourcePrivilegeRepositoryAdminConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", privilege.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", privilege.Description),
					resource.TestCheckResourceAttr(dataSourceName, "repository", privilege.Repository),
					resource.TestCheckResourceAttr(dataSourceName, "format", privilege.Format),
					resource.TestCheckResourceAttr(dataSourceName, "actions.0", string(privilege.Actions[0])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.1", string(privilege.Actions[1])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.2", string(privilege.Actions[2])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.3", string(privilege.Actions[3])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.4", string(privilege.Actions[4])),
				),
			},
		},
	})
}

func testAccDataSourcePrivilegeRepositoryAdminConfig() string {
	return `
	data "nexus_privilege_repository_admin" "acceptance" {
		name = nexus_privilege_repository_admin.acceptance.name
	}`
}
