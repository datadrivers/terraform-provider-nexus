package security_test

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePrivilegeScript(t *testing.T) {
	dataSourceName := "data.nexus_privilege_script.acceptance"

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
				Config: testAccResourceSecurityPrivilegeScriptConfig(privilege) + testAccDataSourcePrivilegeScriptConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", privilege.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", privilege.Description),
					resource.TestCheckResourceAttr(dataSourceName, "script_name", privilege.ScriptName),
					resource.TestCheckResourceAttr(dataSourceName, "actions.0", string(privilege.Actions[0])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.1", string(privilege.Actions[1])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.2", string(privilege.Actions[2])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.3", string(privilege.Actions[3])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.4", string(privilege.Actions[4])),
					resource.TestCheckResourceAttr(dataSourceName, "actions.5", string(privilege.Actions[5])),
				),
			},
		},
	})
}

func testAccDataSourcePrivilegeScriptConfig() string {
	return `
	data "nexus_privilege_script" "acceptance" {
		name = nexus_privilege_script.acceptance.name
	}`
}
