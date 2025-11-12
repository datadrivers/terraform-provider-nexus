package other_test

import (
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
	return `
data "nexus_routing_rule" "acceptance" {
	name = nexus_routing_rule.acceptance.name
}
`
}
