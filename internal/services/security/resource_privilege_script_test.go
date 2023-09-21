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

func TestAccResourceSecurityPrivilegeScript(t *testing.T) {
	resName := "nexus_privilege_script.acceptance"

	privilege := security.PrivilegeScript{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Actions:     []security.SecurityPrivilegeScriptActions{"ADD", "READ", "DELETE", "RUN", "BROWSE", "EDIT"},
		ScriptName:  acctest.RandString(20),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeScriptConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", privilege.Name),
					resource.TestCheckResourceAttr(resName, "description", privilege.Description),
					resource.TestCheckResourceAttr(resName, "script_name", privilege.ScriptName),
					resource.TestCheckResourceAttr(resName, "actions.0", string(privilege.Actions[0])),
					resource.TestCheckResourceAttr(resName, "actions.1", string(privilege.Actions[1])),
					resource.TestCheckResourceAttr(resName, "actions.2", string(privilege.Actions[2])),
					resource.TestCheckResourceAttr(resName, "actions.3", string(privilege.Actions[3])),
					resource.TestCheckResourceAttr(resName, "actions.4", string(privilege.Actions[4])),
					resource.TestCheckResourceAttr(resName, "actions.5", string(privilege.Actions[5])),
				),
			},
		},
	})
}

func testAccResourceSecurityPrivilegeScriptConfig(priv security.PrivilegeScript) string {

	return fmt.Sprintf(`
resource "nexus_script" "acceptance" {
	name = "%s"
	content = "log.info('Hello, World!')"
}

resource "nexus_privilege_script" "acceptance" {
	name = "%s"
	description = "%s"
	actions = [ %s ]
	script_name = resource.nexus_script.acceptance.name
}
`, priv.ScriptName, priv.Name, priv.Description, tools.FormatPrivilegeActionsForConfig(priv.Actions))
}
