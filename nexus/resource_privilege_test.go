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

func TestAccResourcePrivilegeTypeApplication(t *testing.T) {
	var privilege nexus.Privilege

	resName := "nexus_privilege.application"
	priv := nexus.Privilege{
		Actions:     []string{"READ"},
		Description: acctest.RandString(30),
		Domain:      "users",
		Name:        acctest.RandString(10),
		Type:        "application",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivilegeTypeApplicationConfig(priv),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "actions.#", strconv.Itoa(len(priv.Actions))),
					//resource.TestCheckResourceAttr(resName, "actions.0", privilegeActions[0]),
					resource.TestCheckResourceAttr(resName, "description", priv.Description),
					resource.TestCheckResourceAttr(resName, "domain", priv.Domain),
					resource.TestCheckResourceAttr(resName, "name", priv.Name),
					resource.TestCheckResourceAttr(resName, "type", priv.Type),
					testAccCheckPrivilegeResourceExists(resName, &privilege),
				),
			},
			{
				ResourceName:      "nexus_privilege.acceptance",
				ImportStateId:     priv.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourcePrivilegeTypeApplicationConfig(priv nexus.Privilege) string {
	return fmt.Sprintf(`
	resource "nexus_privilege" "application" {
		actions     = ["%s",]
		name        = "%s"
		description = "%s"
		domain      = "%s"
		type        = "%s"
	}
	`, strings.Join(priv.Actions, ",\n"), priv.Name, priv.Description, priv.Domain, priv.Type)
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

func TestAccResourcePrivilegeTypeRepositoryView(t *testing.T) {
	var privilege nexus.Privilege

	resName := "nexus_privilege.repository_view"
	priv := nexus.Privilege{
		Actions:     []string{"READ"},
		Description: acctest.RandString(30),
		Format:      nexus.RepositoryFormatMaven2,
		Name:        acctest.RandString(10),
		Repository:  "maven-releases",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic content selector
			{
				Config: testAccResourcePrivilegeTypeRepositoryViewConfig(priv),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "actions.#", strconv.Itoa(len(priv.Actions))),
					//resource.TestCheckResourceAttr(resName, "actions.0", privilegeActions[0]),
					resource.TestCheckResourceAttr(resName, "description", priv.Description),
					resource.TestCheckResourceAttr(resName, "format", priv.Format),
					resource.TestCheckResourceAttr(resName, "name", priv.Name),
					resource.TestCheckResourceAttr(resName, "type", "repository-view"),
					testAccCheckPrivilegeResourceExists(resName, &privilege),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     priv.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourcePrivilegeTypeRepositoryViewConfig(priv nexus.Privilege) string {
	return fmt.Sprintf(`
resource "nexus_privilege" "repository_view" {
  actions     = ["%s",]
  description = "%s"
  format      = "%s"
  name        = "%s"
  repository  = "%s"
  type        = "repository-view"
}
`, strings.Join(priv.Actions, ",\n"), priv.Description, priv.Format, priv.Name, priv.Repository)
}
