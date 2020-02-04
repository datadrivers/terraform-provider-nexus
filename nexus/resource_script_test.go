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

func testScriptResource(scriptName string, scriptContent string, scriptType string) string {
	return fmt.Sprintf(`
resource "nexus_script" "test" {
    name    = "%s"
	content = "%s"
	type    = "%s"
}
`, scriptName, scriptContent, scriptType))
}

func testAccCheckScriptResourceExists(name string, script *nexus.Script) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		nexusClient := testAccProvider.Meta().(nexus.Client)
		result, err := nexusClient.ScriptRead(rs.Primary.ID)
		if err != nil {
			return err
		}

		*script = *result

		return nil
	}
}

func TestAccScript(t *testing.T) {
	var script nexus.Script

	scriptName := acctest.RandString(10)
	scriptContent := "log.info('Hello, World!')"
	scriptType := "groovy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testScriptResource(scriptName, scriptContent, scriptType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScriptResourceExists("nexus_script.test", &script),
				),
			},
		},
	})
}
