package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestDataSourcePrivileges(t *testing.T) {
	dataSourceName := "data.nexus_privileges.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivileges(),
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
		},
	})
}

func TestDataSourcePrivilegesTypeWildcard(t *testing.T) {
	dataSourceName := "data.nexus_privileges.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivilegesByType(nexus.PrivilegeTypeWildcard),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "id"),
						resource.TestCheckResourceAttr(dataSourceName, "domain", ""),
						resource.TestCheckResourceAttr(dataSourceName, "name", ""),
						resource.TestCheckResourceAttr(dataSourceName, "repository", ""),
						resource.TestCheckResourceAttr(dataSourceName, "type", nexus.PrivilegeTypeWildcard),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "privileges.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "privileges.0.name", "nx-all"),
						resource.TestCheckResourceAttr(dataSourceName, "privileges.0.type", nexus.PrivilegeTypeWildcard),
					),
				),
			},
		},
	})

}

func testAccDataSourcePrivileges() string {
	return fmt.Sprintf(`
data "nexus_privileges" "acceptance" {
}`)
}

func testAccDataSourcePrivilegesByType(privType string) string {
	return fmt.Sprintf(`
data "nexus_privileges" "acceptance" {
	type = "%s"
}`, privType)
}
