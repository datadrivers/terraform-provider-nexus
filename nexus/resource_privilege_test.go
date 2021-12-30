package nexus

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePrivilegeTypeApplication(t *testing.T) {
	var privilege security.Privilege

	resName := "nexus_privilege.application"
	priv := security.Privilege{
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
				ResourceName:      resName,
				ImportStateId:     priv.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourcePrivilegeTypeApplicationConfig(priv security.Privilege) string {
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

func testAccCheckPrivilegeResourceExists(name string, privilege *security.Privilege) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client := testAccProvider.Meta().(*nexus.NexusClient)
		result, err := client.Security.Privilege.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		*privilege = *result

		return nil
	}
}

func TestAccResourcePrivilegeTypeRepositoryView(t *testing.T) {
	var privilege security.Privilege

	resName := "nexus_privilege.repository_view"
	priv := security.Privilege{
		Actions:     []string{"READ"},
		Description: acctest.RandString(30),
		Format:      repository.RepositoryFormatMaven2,
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

func testAccResourcePrivilegeTypeRepositoryViewConfig(priv security.Privilege) string {
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

func TestAccResourcePrivilegeTypeScript(t *testing.T) {

	privilege := security.Privilege{
		Actions:     []string{"READ"},
		Description: acctest.RandString(30),
		Name:        acctest.RandString(10),
		ScriptName:  fmt.Sprintf("sample-script-%s", acctest.RandString(5)),
		Type:        "script",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivilegeResourceTypeScriptConfig(privilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_privilege.script-privilege", "actions.#", strconv.Itoa(len(privilege.Actions))),
					resource.TestCheckResourceAttr("nexus_privilege.script-privilege", "description", privilege.Description),
					resource.TestCheckResourceAttr("nexus_privilege.script-privilege", "name", privilege.Name),
					resource.TestCheckResourceAttr("nexus_privilege.script-privilege", "script_name", privilege.ScriptName),
					resource.TestCheckResourceAttr("nexus_privilege.script-privilege", "type", privilege.Type),
					testAccCheckPrivilegeResourceExists("nexus_privilege.script-privilege", &privilege),
				),
			},
			{
				ResourceName:      "nexus_privilege.script-privilege",
				ImportStateId:     privilege.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourcePrivilegeResourceTypeScriptConfig(priv security.Privilege) string {
	return fmt.Sprintf(`
resource "nexus_script" "%[1]s" {
  name    = "%[1]s"
  content = "log.info('Privilege test')"
  type    = "groovy"
}
resource "nexus_privilege" "script-privilege" {
  actions     = ["%[2]s",]
  description = "%[3]s"
  name        = "%[4]s"
  script_name = "%[1]s"
  type        = "script"
  depends_on  =  ["nexus_script.%[1]s"]
}
`, priv.ScriptName, strings.Join(priv.Actions, ",\n"), priv.Description, priv.Name)
}
