package other_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRoutingRule(t *testing.T) {
	resName := "nexus_routing_rule.acceptance"

	rule := schema.RoutingRule{
		Name:        acctest.RandString(10),
		Description: "acceptance test",
		Mode:        "BLOCK",
		Matchers: []string{
			"^/com/example/.*",
			"^/org/example/.*",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoutingRuleConfig(rule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", rule.Name),
					resource.TestCheckResourceAttr(resName, "name", rule.Name),
					resource.TestCheckResourceAttr(resName, "description", rule.Description),
					resource.TestCheckResourceAttr(resName, "mode", string(rule.Mode)),
					resource.TestCheckResourceAttr(resName, "matchers.#", strconv.Itoa(len(rule.Matchers))),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateId:     rule.Name,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceRoutingRuleConfig(rule schema.RoutingRule) string {
	return fmt.Sprintf(`
	resource "nexus_routing_rule" "acceptance" {
		name        = "%s"
		description = "%s"
		mode        = "%s"
		matchers    = ["%s"]
	  }
`, rule.Name, rule.Description, rule.Mode, strings.Join(rule.Matchers, "\",\""))
}
