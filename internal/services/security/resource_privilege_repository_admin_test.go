package security_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityPrivilegeRepositoryAdmin(t *testing.T) {
	resName := "nexus_privilege_repository_admin.acceptance"

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
				Config: testAccResourceSecurityPrivilegeRepositoryAdminConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", privilege.Name),
					resource.TestCheckResourceAttr(resName, "description", privilege.Description),
					resource.TestCheckResourceAttr(resName, "repository", privilege.Repository),
					resource.TestCheckResourceAttr(resName, "format", privilege.Format),
					//resource.TestCheckResourceAttr(resName, "actions", string(privilege.Actions[])), TODO verify actions
				),
			},
		},
	})
}

func testAccResourceSecurityPrivilegeRepositoryAdminConfig(priv security.PrivilegeRepositoryAdmin) string {

	return fmt.Sprintf(`
	resource "nexus_repository_helm_hosted" "acceptance" {
		name = "%s"
		online = true
	  
		storage {
		  blob_store_name                = "default"
		  strict_content_type_validation = false
		  write_policy                   = "ALLOW"
		}
	}

	resource "nexus_privilege_repository_admin" "acceptance" {
		name = "%s"
		description = "%s"
		actions = [ %s ]
		repository = resource.nexus_repository_helm_hosted.acceptance.name
		format = "%s"
	}
`, priv.Repository, priv.Name, priv.Description, tools.FormatPrivilegeActionsForConfig(priv.Actions), priv.Format)
}
