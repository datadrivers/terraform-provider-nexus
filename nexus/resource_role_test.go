package nexus

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRoleBasic(t *testing.T) {
	t.Parallel()

	roleID := acctest.RandString(10)
	roleName := acctest.RandString(10)
	roleDescription := acctest.RandString(30)
	rolePrivileges := []string{"nx-all"}
	roleRoles := []string{"nx-admin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Creates a basic role
			{
				Config: testAccRoleResource(roleID, roleName, roleDescription, rolePrivileges, roleRoles),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("nexus_role.acceptance", "id", roleID),
					resource.TestCheckResourceAttr("nexus_role.acceptance", "name", roleName),
					resource.TestCheckResourceAttr("nexus_role.acceptance", "roleid", roleID),
					resource.TestCheckResourceAttr("nexus_role.acceptance", "description", roleDescription),
					resource.TestCheckResourceAttr("nexus_role.acceptance", "privileges.#", strconv.Itoa(len(rolePrivileges))),
					resource.TestCheckResourceAttr("nexus_role.acceptance", "roles.#", strconv.Itoa(len(roleRoles))),
				),
			},
			{
				ResourceName:      "nexus_role.acceptance",
				ImportStateId:     roleID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRoleResource(id string, name string, description string, privileges []string, roles []string) string {
	return fmt.Sprintf(`
resource "nexus_role" "acceptance" {
	roleid = "%s"
	name   = "%s"
	description = "%s"
	privileges = ["%s"]
	roles = ["%s"]
}
`, id, name, description, strings.Join(privileges, "\",\""), strings.Join(roles, "\",\""))
}
