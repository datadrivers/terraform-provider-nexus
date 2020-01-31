package nexus

import (
	"fmt"
	"strings"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRole(t *testing.T) {
	var role nexus.Role

	roleID := acctest.RandString(10)
	roleName := acctest.RandString(10)
	roleDescription := acctest.RandString(30)
	rolePrivileges := []string{"nx-all"}
	roleRoles := []string{"nx-admin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResource(roleID, roleName, roleDescription, rolePrivileges, roleRoles),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleResourceExists("nexus_role.test", &role),
				),
			},
			// {
			// 	ResourceName:      "nexus_role.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccRoleResource(id string, name string, description string, privileges []string, roles []string) string {
	return fmt.Sprintf(`
resource "nexus_role" "test" {
	roleid = "%s"
	name   = "%s"
	description = "%s"
	privileges = ["%s"]
	roles = ["%s"]
}
`, id, name, description, strings.Join(privileges, "\",\""), strings.Join(roles, "\",\""))
}

func testAccCheckRoleResourceExists(name string, role *nexus.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		nexusClient := testAccProvider.Meta().(nexus.Client)
		result, err := nexusClient.RoleRead(rs.Primary.ID)
		if err != nil {
			return err
		}

		*role = *result

		return nil
	}
}
