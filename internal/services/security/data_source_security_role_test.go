package security_test

import (
	"strconv"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityRole(t *testing.T) {
	dataSourceName := "data.nexus_security_role.acceptance"

	role := security.Role{
		ID:          acctest.RandString(10),
		Name:        acctest.RandString(10),
		Description: acctest.RandString(30),
		Privileges:  []string{"nx-all"},
		Roles:       []string{"nx-admin"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityRoleConfig(role) + testAccDataSourceSecurityRoleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", role.ID),
					resource.TestCheckResourceAttr(dataSourceName, "roleid", role.ID),
					resource.TestCheckResourceAttr(dataSourceName, "name", role.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", role.Description),
					resource.TestCheckResourceAttr(dataSourceName, "privileges.#", strconv.Itoa(len(role.Privileges))),
					resource.TestCheckResourceAttr(dataSourceName, "roles.#", strconv.Itoa(len(role.Roles))),
				),
			},
		},
	})
}

func testAccDataSourceSecurityRoleConfig() string {
	return `
data "nexus_security_role" "acceptance" {
	roleid = nexus_security_role.acceptance.roleid
}
`
}
