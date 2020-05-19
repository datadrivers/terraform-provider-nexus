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

func TestAccPrivilege(t *testing.T) {
	t.Parallel()

	var privilege nexus.Privilege

	privilegeActions := []string{"READ"}
	privilegeDescription := acctest.RandString(30)
	privilegeDomain := "users"
	privilegeName := acctest.RandString(10)
	privilegeType := "application"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic content selector
			{
				Config: testAccPrivilegeResource(privilegeActions, privilegeName, privilegeDescription, privilegeDomain, privilegeType),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttr("nexus_privilege.acceptance", "actions", privilegeActions),
					resource.TestCheckResourceAttr("nexus_privilege.acceptance", "description", privilegeDescription),
					resource.TestCheckResourceAttr("nexus_privilege.acceptance", "domain", privilegeDomain),
					resource.TestCheckResourceAttr("nexus_privilege.acceptance", "name", privilegeName),
					resource.TestCheckResourceAttr("nexus_privilege.acceptance", "type", privilegeType),
					testAccCheckPrivilegeResourceExists("nexus_privilege.acceptance", &privilege),
				),
			},
			{
				ResourceName:      "nexus_privilege.acceptance",
				ImportStateId:     privilegeName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivilegeResource(actions []string, name string, description string, domain string, tipe string) string {
	return fmt.Sprintf(`
	resource "nexus_privilege" "acceptance" {
		actions = [
			"%s",
		]
		name   = "%s"
		description = "%s"
		domain = "%s"
		type = "%s"
	}
	`, strings.Join(actions, ",\n"), name, description, domain, tipe)
}

func testAccCheckPrivilegeResourceExists(name string, privilege *nexus.Privilege) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client := testAccProvider.Meta().(nexus.Client)
		result, err := client.PrivilegeRead(rs.Primary.ID)
		if err != nil {
			return err
		}

		*privilege = *result

		return nil
	}
}
