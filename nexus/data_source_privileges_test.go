package nexus

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcePrivileges(t *testing.T) {
	dataSourceName := "data.nexus_privileges.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivilegesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "id"),
						resource.TestCheckResourceAttr(dataSourceName, "domain", ""),
						resource.TestCheckResourceAttr(dataSourceName, "name", ""),
						resource.TestCheckResourceAttrSet(dataSourceName, "privileges.#"),
						resource.TestCheckResourceAttrSet(dataSourceName, "privileges.0.name"),
						resource.TestCheckResourceAttr(dataSourceName, "repository", ""),
						resource.TestCheckResourceAttr(dataSourceName, "type", ""),
					),
				),
			},
			{
				Config: testAccDataSourcePrivilegesByTypeConfig(security.PrivilegeTypeWildcard),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "id"),
						resource.TestCheckResourceAttr(dataSourceName, "domain", ""),
						resource.TestCheckResourceAttr(dataSourceName, "name", ""),
						resource.TestCheckResourceAttr(dataSourceName, "repository", ""),
						resource.TestCheckResourceAttr(dataSourceName, "type", security.PrivilegeTypeWildcard),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "privileges.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "privileges.0.name", "nx-all"),
						resource.TestCheckResourceAttr(dataSourceName, "privileges.0.type", security.PrivilegeTypeWildcard),
					),
				),
			},
		},
	})
}

func testAccDataSourcePrivilegesConfig() string {
	return fmt.Sprintf(`
data "nexus_privileges" "acceptance" {
}`)
}

func testAccDataSourcePrivilegesByTypeConfig(privType string) string {
	return fmt.Sprintf(`
data "nexus_privileges" "acceptance" {
	type = "%s"
}`, privType)
}
