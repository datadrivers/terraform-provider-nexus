package nexus

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPrivilegeTypeApplication(t *testing.T) {
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
			{
				Config: testAccPrivilegeResourceTypeApplication(privilegeActions, privilegeName, privilegeDescription, privilegeDomain, privilegeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_privilege.application", "actions.#", strconv.Itoa(len(privilegeActions))),
					//resource.TestCheckResourceAttr("nexus_privilege.application", "actions.0", privilegeActions[0]),
					resource.TestCheckResourceAttr("nexus_privilege.application", "description", privilegeDescription),
					resource.TestCheckResourceAttr("nexus_privilege.application", "domain", privilegeDomain),
					resource.TestCheckResourceAttr("nexus_privilege.application", "name", privilegeName),
					resource.TestCheckResourceAttr("nexus_privilege.application", "type", privilegeType),
					testAccCheckPrivilegeResourceExists("nexus_privilege.application", &privilege),
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

func testAccPrivilegeResourceTypeApplication(actions []string, name string, description string, domain string, tipe string) string {
	return fmt.Sprintf(`
	resource "nexus_privilege" "application" {
		actions     = ["%s",]
		name        = "%s"
		description = "%s"
		domain      = "%s"
		type        = "%s"
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

func TestAccPrivilegeTypeRepositoryView(t *testing.T) {
	var privilege nexus.Privilege

	privilegeActions := []string{"READ"}
	privilegeDescription := acctest.RandString(30)
	privilegeFormat := nexus.RepositoryFormatMaven2
	privilegeName := acctest.RandString(10)
	privilegeRepository := "maven-releases"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic content selector
			{
				Config: testAccPrivilegeResourceTypeRepositoryView(privilegeActions, privilegeDescription, privilegeFormat, privilegeName, privilegeRepository),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_privilege.repository_view", "actions.#", strconv.Itoa(len(privilegeActions))),
					//resource.TestCheckResourceAttr("nexus_privilege.repository_view", "actions.0", privilegeActions[0]),
					resource.TestCheckResourceAttr("nexus_privilege.repository_view", "description", privilegeDescription),
					resource.TestCheckResourceAttr("nexus_privilege.repository_view", "format", privilegeFormat),
					resource.TestCheckResourceAttr("nexus_privilege.repository_view", "name", privilegeName),
					resource.TestCheckResourceAttr("nexus_privilege.repository_view", "type", "repository-view"),
					testAccCheckPrivilegeResourceExists("nexus_privilege.repository_view", &privilege),
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

func testAccPrivilegeResourceTypeRepositoryView(actions []string, description, format, name, repository string) string {
	return fmt.Sprintf(`
resource "nexus_privilege" "repository_view" {
  actions     = ["%s",]
  description = "%s"
  format      = "%s"
  name        = "%s"
  repository  = "%s"
  type        = "repository-view"
}
`, strings.Join(actions, ",\n"), description, format, name, repository)
}
