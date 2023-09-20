package security_test

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePrivilegeRepositoryContentSelector(t *testing.T) {
	dataSourceName := "data.nexus_privilege_repository_content_selector.acceptance"

	privilege := security.PrivilegeRepositoryContentSelector{
		Name:            acctest.RandString(20),
		Description:     acctest.RandString(20),
		Actions:         []security.SecurityPrivilegeRepositoryContentSelectorActions{"ADD", "READ", "DELETE", "BROWSE", "EDIT"},
		Repository:      acctest.RandString(20),
		Format:          "helm",
		ContentSelector: acctest.RandString(20),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeRepositoryContentSelectorConfig(privilege) + testAccDataSourcePrivilegeRepositoryContentSelectorConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", privilege.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", privilege.Description),
					resource.TestCheckResourceAttr(dataSourceName, "content_selector", privilege.ContentSelector),
					resource.TestCheckResourceAttr(dataSourceName, "repository", privilege.Repository),
					resource.TestCheckResourceAttr(dataSourceName, "format", privilege.Format),
					//resource.TestCheckResourceAttr(resName, "actions", string(privilege.Actions[])), TODO verify actions
				),
			},
		},
	})
}

func testAccDataSourcePrivilegeRepositoryContentSelectorConfig() string {
	return `
	data "nexus_privilege_repository_content_selector" "acceptance" {
		name = nexus_privilege_repository_content_selector.acceptance.name
	}`
}
