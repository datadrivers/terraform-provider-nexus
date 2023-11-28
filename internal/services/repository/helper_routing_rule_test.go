package repository_test

import (
	"fmt"
	"strings"

	"github.com/dre2004/go-nexus-client/nexus3/schema"
)

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
