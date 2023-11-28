package security_test

import (
	"fmt"
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityPrivilegeWildcard(t *testing.T) {
	resName := "nexus_privilege_wildcard.acceptance"

	privilege := security.PrivilegeWildcard{
		Name:        acctest.RandString(20),
		Description: acctest.RandString(20),
		Pattern:     "nexus:*",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPrivilegeWildcardConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", privilege.Name),
					resource.TestCheckResourceAttr(resName, "description", privilege.Description),
					resource.TestCheckResourceAttr(resName, "pattern", privilege.Pattern),
				),
			},
		},
	})
}

func testAccResourceSecurityPrivilegeWildcardConfig(priv security.PrivilegeWildcard) string {
	return fmt.Sprintf(`
resource "nexus_privilege_wildcard" "acceptance" {
	name = "%s"
	description = "%s"
	pattern = "%s"
}
`, priv.Name, priv.Description, priv.Pattern)
}
