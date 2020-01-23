package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRole(t *testing.T) {
	t.Parallel()

	roleID := acctest.RandString(10)
	roleName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResource(roleID, roleName),
			},
			{
				ResourceName:      "nexus_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRoleResource(id string, name string) string {
	return fmt.Sprintf(`
resource "nexus_role" "test" {
	roleid = "%s"
	name   = "%s"
}
`, id, name)
}
