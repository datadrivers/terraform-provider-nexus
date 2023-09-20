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

func TestAccResourceSecurityPrivilegeRepositoryContentSelector(t *testing.T) {
	resName := "nexus_privilege_repository_content_selector.acceptance"

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
				Config: testAccResourceSecurityPrivilegeRepositoryContentSelectorConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", privilege.Name),
					resource.TestCheckResourceAttr(resName, "description", privilege.Description),
					resource.TestCheckResourceAttr(resName, "repository", privilege.Repository),
					resource.TestCheckResourceAttr(resName, "format", privilege.Format),
					resource.TestCheckResourceAttr(resName, "content_selector", privilege.ContentSelector),
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

func testAccResourceSecurityPrivilegeRepositoryContentSelectorConfig(priv security.PrivilegeRepositoryContentSelector) string {

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

	resource "nexus_security_content_selector" "acceptance" {
		name        = "%s"
		description = "A content selector matching public docker images."
		expression  = "path =^ \"/v2/public/\""
	}

	resource "nexus_privilege_repository_content_selector" "acceptance" {
		name = "%s"
		description = "%s"
		actions = [ %s ]
		repository = resource.nexus_repository_helm_hosted.acceptance.name
		format = "%s"
		content_selector = resource.nexus_security_content_selector.acceptance.name
	}
`, priv.Repository, priv.ContentSelector, priv.Name, priv.Description, tools.FormatPrivilegeActionsForConfig(priv.Actions), priv.Format)
}
