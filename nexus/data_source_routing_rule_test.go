package nexus

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRoutingRule(t *testing.T) {
	resName := "data.nexus_routing_rule.acceptance"
	rule := schema.RoutingRule{
		Name:        "acceptance",
		Description: "acceptance test",
		Mode:        "BLOCK",
		Matchers: []string{
			"^/com/example/.*",
			"^/org/example/.*",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoutingRuleConfig(rule) + testAccDataSourceRoutingRuleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", rule.Name),
					resource.TestCheckResourceAttr(resName, "description", rule.Description),
					resource.TestCheckResourceAttr(resName, "mode", string(rule.Mode)),
					resource.TestCheckResourceAttrSet(resName, "matchers.#"),
				),
			},
		},
	})
}

func testAccDataSourceRoutingRuleConfig() string {
	return fmt.Sprintf(`
data "nexus_routing_rule" "acceptance" {
	name = nexus_routing_rule.acceptance.name
}
`)
}
