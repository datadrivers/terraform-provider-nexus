package security_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityPrivilegeApplication(t *testing.T) {
	resName := "nexus_privilege_application.acceptance"

	privilege := security.PrivilegeApplication{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Actions:     []security.SecurityPrivilegeApplicationActions{"DELETE"},
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
	actions = ["DELETE"]
	domain = "%s"
}
`, priv.Name, priv.Description, priv.Domain)
}
