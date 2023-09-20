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

func TestAccResourceSecurityPrivilegeRepositoryView(t *testing.T) {
	resName := "nexus_privilege_repository_view.acceptance"

	privilege := security.PrivilegeRepositoryView{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Actions:     []security.SecurityPrivilegeRepositoryViewActions{"ADD", "READ", "DELETE", "BROWSE", "EDIT"},
		Repository:  acctest.RandString(20),
		Format:      "helm",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeRepositoryViewConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", privilege.Name),
					resource.TestCheckResourceAttr(resName, "description", privilege.Description),
					resource.TestCheckResourceAttr(resName, "repository", privilege.Repository),
					resource.TestCheckResourceAttr(resName, "format", privilege.Format),
					resource.TestCheckResourceAttr(resName, "actions.0", string(privilege.Actions[0])),
					resource.TestCheckResourceAttr(resName, "actions.1", string(privilege.Actions[1])),
					resource.TestCheckResourceAttr(resName, "actions.2", string(privilege.Actions[2])),
					resource.TestCheckResourceAttr(resName, "actions.3", string(privilege.Actions[3])),
					resource.TestCheckResourceAttr(resName, "actions.4", string(privilege.Actions[4])),
				),
			},
		},
	})
}

func testAccResourceSecurityPrivilegeRepositoryViewConfig(priv security.PrivilegeRepositoryView) string {
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

resource "nexus_privilege_repository_view" "acceptance" {
	name = "%s"
	description = "%s"
	actions = [ %s ]
	repository = resource.nexus_repository_helm_hosted.acceptance.name
	format = "%s"
}
`, priv.Repository, priv.Name, priv.Description, tools.FormatPrivilegeActionsForConfig(priv.Actions), priv.Format)
}
