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

func TestAccResourceSecurityPrivilegeApplication(t *testing.T) {
	resName := "nexus_privilege_application.acceptance"

	privilege := security.PrivilegeApplication{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Actions:     []security.SecurityPrivilegeApplicationActions{"ADD", "READ", "EDIT", "DELETE", "STOP", "DISASSOCIATE", "ASSOCIATE", "START", "ALL"},
		Domain:      acctest.RandString(20),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeApplicationConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", privilege.Name),
					resource.TestCheckResourceAttr(resName, "description", privilege.Description),
					resource.TestCheckResourceAttr(resName, "domain", privilege.Domain),
					resource.TestCheckResourceAttr(resName, "actions.0", string(privilege.Actions[0])),
					resource.TestCheckResourceAttr(resName, "actions.1", string(privilege.Actions[1])),
					resource.TestCheckResourceAttr(resName, "actions.2", string(privilege.Actions[2])),
					resource.TestCheckResourceAttr(resName, "actions.3", string(privilege.Actions[3])),
					resource.TestCheckResourceAttr(resName, "actions.4", string(privilege.Actions[4])),
					resource.TestCheckResourceAttr(resName, "actions.5", string(privilege.Actions[5])),
					resource.TestCheckResourceAttr(resName, "actions.6", string(privilege.Actions[6])),
					resource.TestCheckResourceAttr(resName, "actions.7", string(privilege.Actions[7])),
					resource.TestCheckResourceAttr(resName, "actions.8", string(privilege.Actions[8])),
				),
			},
		},
	})
}

func testAccResourceSecurityPrivilegeApplicationConfig(priv security.PrivilegeApplication) string {
	return fmt.Sprintf(`
resource "nexus_privilege_application" "acceptance" {
	name = "%s"
	description = "%s"
	actions = [ %s ]
	domain = "%s"
}
`, priv.Name, priv.Description, tools.FormatPrivilegeActionsForConfig(priv.Actions), priv.Domain)
}
